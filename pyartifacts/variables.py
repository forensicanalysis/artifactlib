# Copyright (c) 2019 Siemens AG
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#
# Author(s): Demian Kellermann

# pylint: disable=too-few-public-methods

""" Module for helpers regarding variable expansion """

import logging
import re
from collections import defaultdict
from typing import Set, Iterable, Dict, Callable, TYPE_CHECKING

import networkx

from .definitions import SOURCE_TYPE_DIRECTORY, SOURCE_TYPE_FILE, SOURCE_TYPE_PATH, SOURCE_TYPE_REGISTRY_KEY, \
    SOURCE_TYPE_REGISTRY_VALUE

if TYPE_CHECKING:  # we need this disabled at runtime since it's a circular import
    from .artifact import ArtifactDefinition, ArtifactSource # pylint: disable=cyclic-import

LOGGER = logging.getLogger(__name__)
VARIABLE_IDENTIFIER = re.compile(r'(%?%([a-zA-Z0-9_.-]+)%?%)')


def get_needed_vars(*entries: str) -> Iterable[str]:
    """ Parses all used variables from a string """
    results = set()
    for entry in entries:
        variables = VARIABLE_IDENTIFIER.findall(entry)  # will return tuples: (with surrounding %s, without)
        results.update([match[1] for match in variables])
    return results


class KnowledgeBase:
    """ Class to collect resolved variables and manage their dependencies """

    def __init__(self, artifacts: Dict[str, 'ArtifactDefinition']):
        self.artifacts = artifacts
        self.providers = defaultdict(set)
        self.graph = networkx.DiGraph()
        self.resolved_vars: Dict[str, Iterable[str]] = {}
        self._build_graph(artifacts)

    def get_value(self, key: str, resolve_callback: Callable[['ArtifactSource'], Iterable[str]]) -> Iterable[str]:
        """
        Retrieve a value for a variable. In case the value is not yet available,
        this will trigger  evaluation of the variable.
        To facilitate this, a callback must be provided that can be used to resolve
        another ArtifactSource to (text) values. The text value is expected to be
            - For REGISTRY_KEY, the key path(s)
            - For REGISTRY_VALUE, the string representation of the value(s) indicated
            - For PATH and DIRECTORY, the file system path(s)
            - For FILE, the content of the file(s) in one '\n'-separated string per source
        In any case, the value(s) will be cached for further use before returning to the caller.
        """
        value = self.resolved_vars.get(key, None)
        if value is not None:
            return value

        # get the providers of this value (in-edges in the graph)
        if key not in self.graph:
            raise ValueError("No providers found for %s" % key)

        collected_values: Set[str] = set()
        providers = self.graph.in_edges(nbunch=[key])
        if not providers:
            raise ValueError("No providers are registered for %s" % key)
        for provider, __ in providers:
            provider_result = resolve_callback(provider)
            collected_values.update(self._extract_var(key, provider, provider_result))

        self.resolved_vars[key] = collected_values
        return collected_values

    @staticmethod
    def _extract_var(key: str, source: 'ArtifactSource', data: Iterable[str]) -> Iterable[str]:
        """ Extracts variable matches from a list of results, taking regexes and different types into account """
        if source.type == SOURCE_TYPE_FILE:
            all_data = [line.strip() for text in data for line in text.split('\n')]
        elif source.type in (SOURCE_TYPE_PATH, SOURCE_TYPE_DIRECTORY,
                             SOURCE_TYPE_REGISTRY_KEY, SOURCE_TYPE_REGISTRY_VALUE):
            all_data = data
        else:
            raise ValueError("Unsupported source type for variable expansion: %s" % source.type)

        for provider in source.provides:
            if provider.key == key:  # have to fish out the right provides-directive
                if provider.regex:
                    return [m.group(1) for l in all_data for m in (re.search(provider.regex, l),) if m]
                return all_data
        raise ValueError("No provider matched %s" % key)

    def _build_graph(self, artifacts: Dict[str, 'ArtifactDefinition']) -> None:
        for artifact in artifacts.values():
            for source in artifact.sources:
                for provide in source.provides:
                    self.providers[provide.key].add(source)
                    # If an artifact provides some variable, make an incoming edge to this variable
                    self.graph.add_node(source, type='source')
                    self.graph.add_node(provide.key, type='variable')
                    self.graph.add_edge(source, provide.key)
                for variable in source.needs:
                    # for every dependency to a variable, make an outgoing edge from that variable
                    self.graph.add_node(source, type='source')
                    self.graph.add_node(variable, type='variable')
                    self.graph.add_edge(variable, source)

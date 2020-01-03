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

""" The registry keeps track of all known artifacts """

import json
import logging
import os
from typing import Dict, Optional, Iterable

import jsonschema
import yaml

from .artifact import ArtifactDefinition, make_artifact
from .variables import KnowledgeBase
from .definitions import OS_TYPES

LOGGER = logging.getLogger(__name__)


class Registry:

    def __init__(self):
        self.artifacts: Dict[str, ArtifactDefinition] = {}
        self.artifact_schema = None

        self._load_schema()

    def read_folder(self, folder: str, os_filter: Optional[Iterable[str]] = None) -> None:
        """
        Reads all artifact definition files in a local folder and adds them to the registry.
        Search is not recursive
        """
        input_files = [f for f in os.listdir(folder) if f.endswith('.yaml')]
        LOGGER.debug("Reading files from %s: %s", folder, input_files)
        for input_file in input_files:
            self.read_file(os.path.join(folder, input_file), os_filter)

    def read_file(self, file: str, os_filter: Optional[Iterable[str]] = None) -> None:
        """ Reads a single YAML file and adds the contained artifacts to the registry """
        if not os_filter:
            os_filter = OS_TYPES
        with open(file, 'r') as definition_file:
            try:
                artifacts = list(yaml.safe_load_all(definition_file))
            except yaml.YAMLError as err:
                LOGGER.error("Could not load %s: %s", file, err)
                return
        for artifact_dict in artifacts:
            if self.verify_structure(artifact_dict):
                artifact = make_artifact(artifact_dict)
                if artifact:
                    if artifact.supported_os and not any(os_name in os_filter for os_name in artifact.supported_os):
                        continue
                    self.artifacts[artifact.name] = artifact

            else:
                LOGGER.error("Not adding %s due to validation error", artifact_dict.get('name', '!ERR_NO_NAME'))

    def get_knowledge_base(self) -> KnowledgeBase:
        return KnowledgeBase(self.artifacts)

    def verify_structure(self, artifact_dict: dict) -> bool:
        """ Checks if an artifact definition is valid can be parsed into a definition object """
        try:
            jsonschema.validate(artifact_dict, self.artifact_schema)
            return True
        except jsonschema.exceptions.ValidationError as err:
            LOGGER.error("Artifact is not valid: %s", err)
            return False

    def get_artifact(self, title: str) -> Optional[ArtifactDefinition]:
        """ Return the artifact with the supplied name, if present in the registry """
        return self.artifacts.get(title, None)

    def _load_schema(self):
        my_path = os.path.dirname(os.path.realpath(__file__))
        my_schema = os.path.join(my_path, 'artifact_schema.json')
        with open(my_schema, 'r') as schema:
            self.artifact_schema = json.load(schema)

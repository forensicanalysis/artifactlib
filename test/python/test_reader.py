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

import os
import logging

import pytest

import pyartifacts

TEST_ARTIFACT_LOCATION = os.path.join(os.path.dirname(__file__), '..', 'artifacts', "test.yaml")
ARTIFACT_LOCATION = os.environ.get('ARTIFACT_PATH', TEST_ARTIFACT_LOCATION)

logging.basicConfig(level=os.environ.get('LOGLEVEL', 'DEBUG'))


@pytest.fixture
def registry():
    registry = pyartifacts.Registry()
    registry.read_file(ARTIFACT_LOCATION)
    return registry


class TestPyArtifacts:

    def test_artifact_read(self, registry):
        assert len(registry.artifacts) == 7
        assert 'EmptyArtifact' not in registry.artifacts


    def test_variable_resolving(self, registry):
        replies = {
            'HKEY_LOCAL_MACHINE\\Software\\Microsoft\\Windows NT\\CurrentVersion\\ProfileList\\*':
                [
                    'HKEY_LOCAL_MACHINE\\Software\\Microsoft\\Windows NT\\CurrentVersion\\ProfileList\\S-1337',
                    'HKEY_LOCAL_MACHINE\\Software\\Microsoft\\Windows NT\\CurrentVersion\\ProfileList\\X-1111',
                    'HKEY_LOCAL_MACHINE\\Software\\Microsoft\\Windows NT\\CurrentVersion\\ProfileList\\FOO-S-BAR',
                    'HKEY_LOCAL_MACHINE\\Software\\Microsoft\\Windows NT\\CurrentVersion\\ProfileList\\X-1111',
                ]
        }

        def callback(source):
            if source.type == 'REGISTRY_KEY':
                return replies[source.keys[0]]
            else:
                assert False

        kb: pyartifacts.KnowledgeBase = registry.get_knowledge_base()
        resolved = kb.get_value('users.sid', callback)
        assert resolved == {'S-1337', 'X-1111', 'FOO-S-BAR'}

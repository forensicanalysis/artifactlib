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

""" Constants used for types """

SOURCE_TYPE_ARTIFACT_GROUP = "ARTIFACT_GROUP"
SOURCE_TYPE_COMMAND = "COMMAND"
SOURCE_TYPE_DIRECTORY = "DIRECTORY"
SOURCE_TYPE_FILE = "FILE"
SOURCE_TYPE_PATH = "PATH"
SOURCE_TYPE_REGISTRY_KEY = "REGISTRY_KEY"
SOURCE_TYPE_REGISTRY_VALUE = "REGISTRY_VALUE"
SOURCE_TYPE_WMI = "WMI"

SOURCE_TYPES = [
    SOURCE_TYPE_ARTIFACT_GROUP,
    SOURCE_TYPE_COMMAND,
    SOURCE_TYPE_DIRECTORY,
    SOURCE_TYPE_FILE,
    SOURCE_TYPE_PATH,
    SOURCE_TYPE_REGISTRY_KEY,
    SOURCE_TYPE_REGISTRY_VALUE,
    SOURCE_TYPE_WMI,
]

OS_WINDOWS = 'Windows'
OS_LINUX = 'Linux'
OS_DARWIN = 'Darwin'

OS_TYPES = [
    OS_WINDOWS,
    OS_LINUX,
    OS_DARWIN
]

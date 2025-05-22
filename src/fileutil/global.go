/*
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 Written by Frederic PONT.
 (c) Frederic Pont 2025
*/

package fileutil

type Args struct {
	ComparisonCriteria string // comparison by filename or path (to compare files with same file names)
}

// Define the struct for software parameters
type Conf struct {
	Arg            Args
	TrimIndexPathA int
	TrimIndexPathB int
	TrimPath       bool
	CommonDirSep   string
}

// Declare a global variable of type Conf
var Config Conf

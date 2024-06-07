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
 (c) Frederic Pont 2024
*/

package fileutil

import (
	"fmt"
	"time"
)

func Timer(countdownFrom int) {

	fmt.Printf("This window will close after %d seconds:\n", countdownFrom)
	for i := countdownFrom; i > 0; i-- {
		fmt.Print(i, "\r")
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Countdown finished!")

}

// A C level header file for the NES emulation library.
// Copyright 2021 Christian Kauten
//
// Author: Christian Kauten (kautencreations@gmail.com)
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//

#ifndef NES_EMULATOR_H_
#define NES_EMULATOR_H_

/// Return the width of the NES.
extern int Width();

/// Return the height of the NES.
extern int Height();

/// Initialize a new emulator and return a pointer to it
extern void* Initialize(char*);

/// Return a pointer to a controller on the machine
extern char* Controller(void*, int);

/// Return the pointer to the screen buffer
extern char* Screen(void*);

/// Return the pointer to the memory buffer
extern char* Memory(void*);

/// Reset the emulator
extern void Reset(void*);

/// Perform a discrete step in the emulator (i.e., 1 frame)
extern void Step(void*);

/// Create a deep copy (i.e., a clone) of the given emulator
extern void Backup(void*);

/// Restore from a deep copy (i.e., a clone) of the given emulator
extern void Restore(void*);

/// Close the emulator, i.e., purge it from memory
extern void Close(void*);

#endif  // NES_EMULATOR_H_

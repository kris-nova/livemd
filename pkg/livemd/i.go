/*===========================================================================*\
 *           MIT License Copyright (c) 2022 Kris Nóva <kris@nivenly.com>     *
 *                                                                           *
 *                ┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓                *
 *                ┃   ███╗   ██╗ ██████╗ ██╗   ██╗ █████╗   ┃                *
 *                ┃   ████╗  ██║██╔═████╗██║   ██║██╔══██╗  ┃                *
 *                ┃   ██╔██╗ ██║██║██╔██║██║   ██║███████║  ┃                *
 *                ┃   ██║╚██╗██║████╔╝██║╚██╗ ██╔╝██╔══██║  ┃                *
 *                ┃   ██║ ╚████║╚██████╔╝ ╚████╔╝ ██║  ██║  ┃                *
 *                ┃   ╚═╝  ╚═══╝ ╚═════╝   ╚═══╝  ╚═╝  ╚═╝  ┃                *
 *                ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛                *
 *                                                                           *
 *                       This machine kills fascists.                        *
 *                                                                           *
\*===========================================================================*/

package livemd

import (
	"io/ioutil"
	"strconv"
)

const (
	IPerm = 0655
	IFile = "i"
)

// I will retrieve the local index
func I() int {
	bytes, err := ioutil.ReadFile(IFile)
	if err != nil {
		err = ioutil.WriteFile(IFile, []byte("0"), IPerm)
		if err != nil {
			return -1
		}
		return I()
	}
	i, err := strconv.Atoi(string(bytes))
	if err != nil {
		return -1
	}
	return i
}

// ISet will update the local index
func ISet(x int) int {
	str := strconv.Itoa(x)
	err := ioutil.WriteFile(IFile, []byte(str), IPerm)
	if err != nil {
		return -1
	}
	return x
}

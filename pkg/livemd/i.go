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

	"github.com/sirupsen/logrus"
)

const (
	IPerm = 0655
	IFile = "i"
)

// I will retrieve the local index
func I() int {
	bytes, err := ioutil.ReadFile(IFile)
	if err != nil {
		logrus.Infof("Index file reset: 0")
		err = ioutil.WriteFile(IFile, []byte("0"), IPerm)
		if err != nil {
			logrus.Errorf(err.Error())
			return -1
		}
		return I()
	}
	i, err := strconv.Atoi(string(bytes))
	if err != nil {
		logrus.Infof("Index file reset: 0")
		err = ioutil.WriteFile(IFile, []byte("0"), IPerm)
		if err != nil {
			logrus.Errorf(err.Error())
			return -1
		}
		return I()
	}
	return i
}

// ISet will update the local index
func ISet(x int) int {
	str := strconv.Itoa(x)
	err := ioutil.WriteFile(IFile, []byte(str), IPerm)
	if err != nil {
		logrus.Errorf(err.Error())
		return -1
	}
	return x
}

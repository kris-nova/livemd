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
	"path/filepath"
	"strconv"

	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

const (
	IPerm    = 0655
	defaultI = "i"
)

// i will return the path of the i file
func i() string {
	dir, err := homedir.Dir()
	if err != nil {
		return defaultI
	}
	return filepath.Join(dir, ".livemd", defaultI)
}

// I will retrieve the local index
func I() int {
	bytes, err := ioutil.ReadFile(i())
	if err != nil {
		logrus.Infof("Index file reset: 0")
		err = ioutil.WriteFile(i(), []byte("0"), IPerm)
		if err != nil {
			logrus.Errorf(err.Error())
			return -1
		}
		return I()
	}
	iint, err := strconv.Atoi(string(bytes))
	if err != nil {
		logrus.Infof("Index file reset: 0")
		err = ioutil.WriteFile(i(), []byte("0"), IPerm)
		if err != nil {
			logrus.Errorf(err.Error())
			return -1
		}
		return I()
	}
	return iint
}

// ISet will update the local index
func ISet(x int) int {
	str := strconv.Itoa(x)
	err := ioutil.WriteFile(i(), []byte(str), IPerm)
	if err != nil {
		logrus.Errorf(err.Error())
		return -1
	}
	return x
}

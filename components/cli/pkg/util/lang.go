/*
 * Copyright (c) 2019 WSO2 Inc. (http:www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http:www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package util

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type LangManager interface {
	Init() error
}

type BLangManager struct{}

func (langMgr *BLangManager) Init() error {
	userRepo := filepath.Join(UserHomeDir(), ".ballerina")
	// if not exists, create the location
	if _, err := os.Stat(userRepo); os.IsNotExist(err) {
		if err = os.Mkdir(userRepo, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	// copy from cellery installation location to user repository
	installPath := filepath.Join(CelleryInstallationDir(), "repo")
	cmd := exec.Command("cp", "-r", installPath, userRepo)
	stderrReader, _ := cmd.StderrPipe()
	stderrScanner := bufio.NewScanner(stderrReader)

	execError := ""
	go func() {
		for stderrScanner.Scan() {
			execError += stderrScanner.Text()
		}
	}()
	if err := cmd.Run(); err != nil {
		log.Printf("Error: %v \n", execError)
		return err
	}
	return nil
}

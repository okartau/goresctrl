/*
Copyright 2021 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sst

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type cpuPackageInfo struct {
	id   int
	cpus []int
}

func getOnlineCpuPackages() (map[int]*cpuPackageInfo, error) {
	basePath := "/sys/bus/cpu/devices"

	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	pkgs := make(map[int]*cpuPackageInfo)

	for _, file := range files {
		// Try to read siblings from topology
		raw, err := ioutil.ReadFile(filepath.Join(basePath, file.Name(), "topology/physical_package_id"))
		if os.IsNotExist(err) {
			// Offline -> topology information does not exist
			continue
		} else if err != nil {
			return nil, err
		}

		cpuId, err := strconv.Atoi(file.Name()[3:])
		if err != nil {
			return nil, err
		}

		pkgId, err := strconv.Atoi(strings.TrimSpace(string(raw)))
		if err != nil {
			return nil, err
		}

		if _, ok := pkgs[pkgId]; !ok {
			pkgs[pkgId] = &cpuPackageInfo{id: pkgId}
		}
		pkgs[pkgId].cpus = append(pkgs[pkgId].cpus, cpuId)
	}

	return pkgs, nil
}

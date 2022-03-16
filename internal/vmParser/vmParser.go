package vmParser

import (
	"awesomeProjectRucenter/internal/model"
	"awesomeProjectRucenter/pkg/erx"
	"awesomeProjectRucenter/pkg/tools"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sort"
	"sync"
)

var log *logrus.Logger

func init() {
	log = tools.GetLogrusInstance("")
}

type diskResults struct {
	model.DiskResults
}

type vmResults struct {
	model.VmResults
}

type getter interface {
	GetData(user, pass, url string) error
}

type syncMapId struct {
	sync.Mutex
	vmMap map[int][]model.Disk
}

type syncMapUuid struct {
	sync.Mutex
	vmMap map[string]model.VmDiscs
}

func PrintSyncMap(domain string) (*syncMapId, error) {
	var diskParsed, vmParsed interface{}
	var err error
	host := fmt.Sprintf("http://%v", domain)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		diskParsed, err = parseResults(diskResults{model.DiskResults{Next: host + "/disks/"}})
		if err != nil {
			err = erx.New(err)
		}
	}()
	go func() {
		defer wg.Done()
		vmParsed, err = parseResults(vmResults{model.VmResults{Next: host + "/vms/"}})
		if err != nil {
			err = erx.New(err)
		}
	}()
	wg.Wait()
	fmt.Println("wg done")
	if err != nil {
		return nil, err
	}
	vmRes := vmParsed.(vmResults)
	diskRes := diskParsed.(diskResults)
	sMap := new(syncMapId)
	sMap = writeMapCycle(vmRes, diskRes)

	return sMap, err
}

func PrintAsyncMap(domain string) (*syncMapUuid, error) {
	var diskParsed, vmParsed interface{}
	var err error
	host := fmt.Sprintf("http://%v", domain)
	sMap := new(syncMapUuid)
	sMap.vmMap = make(map[string]model.VmDiscs)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		diskParsed, err = parseResults(diskResults{model.DiskResults{Next: host + "/disks/"}})
		if err != nil {
			err = erx.New(err)
			return
		}
		diskRes := diskParsed.(diskResults)
		asyncWriteMap(diskRes, sMap)
	}()
	go func() {
		defer wg.Done()
		vmParsed, err = parseResults(vmResults{model.VmResults{Next: host + "/vms/"}})
		if err != nil {
			err = erx.New(err)
			return
		}
		vmRes := vmParsed.(vmResults)
		asyncWriteMap(vmRes, sMap)
	}()
	wg.Wait()

	return sMap, err
}

func writeMapCycle(vmRes vmResults, diskRes diskResults) *syncMapId {
	resMap := syncMapId{vmMap: make(map[int][]model.Disk)}
	for _, vm := range vmRes.Results {
		for _, disk := range diskRes.Results {
			if disk.Vm == vm.Uuid {
				resMap.vmMap[vm.Id] = append(resMap.vmMap[vm.Id], disk)
			}
		}
	}
	return &resMap
}

func asyncWriteMap(i interface{}, m *syncMapUuid) {
	m.Lock()
	defer m.Unlock()
	switch v := i.(type) {
	case diskResults:
		for _, v := range v.Results {
			val, ok := m.vmMap[v.Vm]
			if ok {
				val.Discs = append(val.Discs, v)
				m.vmMap[v.Vm] = val
			}
			if !ok {
				m.vmMap[v.Vm] = model.VmDiscs{Discs: []model.Disk{v}}
			}
		}
	case vmResults:
		for _, v := range v.Results {
			val, ok := m.vmMap[v.Uuid]
			if ok && val.Uuid == "" {
				val.Id = v.Id
				val.Uuid = v.Uuid
				val.Name = v.Name
				m.vmMap[v.Uuid] = val
			}
			if !ok {
				m.vmMap[v.Uuid] = model.VmDiscs{Vm: v}
			}
		}
	}
}

func PrintSortedMap(m *syncMapId) {
	var s []int
	for i := range m.vmMap {
		s = append(s, i)
	}
	sort.Ints(s)
	for _, v := range s {
		fmt.Println(m.vmMap[v])
	}

}

func parseResults(g interface{}) (interface{}, error) {
	switch v := g.(type) {
	case diskResults:
		for url := v.Next; url != ""; url = v.Next {
			err := v.GetData("admin", "admin", v.Next)
			if err != nil {
				return nil, err
			}
			if url == v.Next {
				break
			}
		}
		return v, nil
	case vmResults:
		for url := v.Next; url != ""; url = v.Next {
			err := v.GetData("admin", "admin", v.Next)
			if err != nil {
				return nil, err
			}
			if url == v.Next {
				break
			}
		}
		return v, nil
	default:
		return nil, erx.New(fmt.Errorf("type assertion failed"))
	}
	return nil, nil
}

func (r *diskResults) GetData(user, pass, url string) error {
	log.Info("Getting " + url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(user, pass)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(erx.New(err))
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var results diskResults
	err = json.Unmarshal([]byte(bodyText), &results)
	if err != nil {
		return err
	}
	r.Next = results.Next
	r.Previous = results.Previous
	r.Results = append(r.Results, results.Results...)
	log.Info("Success")
	return err
}

func (r *vmResults) GetData(user, pass, url string) error {
	log.Info("Getting " + url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(user, pass)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(erx.New(err))
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var results vmResults
	err = json.Unmarshal([]byte(bodyText), &results)
	if err != nil {
		return err
	}
	r.Next = results.Next
	r.Previous = results.Previous
	r.Results = append(r.Results, results.Results...)
	log.Info("Success")
	return err
}

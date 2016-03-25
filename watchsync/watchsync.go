package watchsync

import (
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"

	"gopkg.in/fsnotify.v1"
)

//WatchSync watch sync
type WatchSync struct {
	Host         string
	User         string
	Port         string
	Src          string
	Dest         string
	ExcludeSync  []string
	ExcludeWatch []string
}

//Watch Monitoring files.
func (s *WatchSync) Watch() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("failed watcher: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write:
					s.Sync()
				case event.Op&fsnotify.Create == fsnotify.Create:
					s.Sync()
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					s.Sync()
				}
			case err := <-watcher.Errors:
				log.Fatal("failed watcher errors: ", err)
				done <- true
			}
		}
	}()

	//監視対象ディレクトリ取得
	dirs, err := s.getDirs(s.Src)
	if err != nil {
		log.Fatal("failed dirs: ", err)
	}

	//監視
	for _, dir := range dirs {
		err = watcher.Add(dir)
		if err != nil {
			log.Fatal("failed watcher add: ", err)
		}
	}

	log.Println("### Watching! ", s.Src, " ###")

	<-done
}

//Sync sync.
func (s *WatchSync) Sync() {
	cmd := "/usr/bin/rsync"

	exclude := []string{}
	for _, e := range s.ExcludeSync {
		exclude = append(exclude, "--exclude="+e)
	}

	args := []string{}
	args = append(args, []string{
		"-avz",
		"--delete",
		"-e",
		"/usr/bin/ssh -p " + s.Port,
	}...)
	args = append(args, exclude...)
	args = append(args, []string{
		s.Src,
		s.User + "@" + s.Host + ":" + s.Dest,
	}...)

	log.Println("### Starting! ###")

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Fatal("failed command: ", err)
	}

	log.Println(string(out))
	log.Println("### Finished! ###")
}

func (s *WatchSync) getDirs(path string) ([]string, error) {

	dirs := []string{}
	dirs = append(dirs, path)

	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}

		newDirPath := filepath.Join(path, fi.Name())

		if s.checkExcludeWatchDir(newDirPath) {
			continue
		}

		newDirs, err := s.getDirs(newDirPath)
		if err != nil {
			return nil, err
		}

		dirs = append(dirs, newDirs...)
	}

	return dirs, nil
}

func (s *WatchSync) checkExcludeWatchDir(path string) bool {
	ret := false

	for _, e := range s.ExcludeWatch {
		if s.Src+"/"+e == path {
			ret = true
			break
		}
	}

	return ret
}

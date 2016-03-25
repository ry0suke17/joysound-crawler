package main

import "bitbucket.org/yneee/exsongs/watchsync"

func main() {
	s := watchsync.WatchSync{
		Host: "8887.vodsys.net",
		User: "admin2",
		Port: "22",
		Dest: "/home/admin2/docker/golang/src/bitbucket.org/yneee/",

		Src: "/Users/y/.gvm/pkgsets/go1.5/global/src/bitbucket.org/yneee/exsongs",
		ExcludeSync: []string{
			".git",
			"node_modules",
			"bower_components",
			".sass-cache",
			".gulp-scss-cache",
			".DS_Store",
			".AppleDouble",
			".LSOverride",
			"vendor",
			".idea",
		},
		ExcludeWatch: []string{
			".git",
			"node_modules",
			"bower_components",
			".sass-cache",
			".gulp-scss-cache",
			".DS_Store",
			".AppleDouble",
			".LSOverride",
			"vendor",
			".idea",
		},
	}

	s.Watch()
}

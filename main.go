package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	longString = []struct {
		key   string
		value string
	}{
		{"emerald", "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Nunc faucibus a pellentesque sit. Mauris commodo quis imperdiet massa tincidunt nunc pulvinar sapien. In dictum non consectetur a erat nam at. Pellentesque habitant morbi tristique senectus et. Tellus integer feugiat scelerisque varius morbi enim nunc faucibus a. Dictumst vestibulum rhoncus est pellentesque elit ullamcorper dignissim. Turpis massa sed elementum tempus egestas sed sed risus. Fringilla urna porttitor rhoncus dolor purus non enim praesent. Aliquet bibendum enim facilisis gravida neque convallis a. Volutpat consequat mauris nunc congue nisi vitae suscipit tellus mauris. At imperdiet dui accumsan sit amet nulla facilisi. Justo nec ultrices dui sapien eget mi proin. Integer quis auctor elit sed vulputate mi sit amet mauris. Consequat ac felis donec et odio pellentesque diam. Eget sit amet tellus cras adipiscing enim. Quis commodo odio aenean sed adipiscing diam donec. Sit amet luctus venenatis lectus magna fringilla."},
		{"emerald", "Montes nascetur ridiculus mus mauris vitae. Nisi est sit amet facilisis. Amet est placerat in egestas erat imperdiet sed euismod nisi. Sagittis vitae et leo duis ut diam. Hac habitasse platea dictumst vestibulum rhoncus est. A arcu cursus vitae congue. Tortor pretium viverra suspendisse potenti. Mauris pharetra et ultrices neque ornare aenean euismod. Blandit turpis cursus in hac habitasse platea. Sed nisi lacus sed viverra tellus in hac habitasse. Euismod elementum nisi quis eleifend quam adipiscing vitae. At ultrices mi tempus imperdiet nulla malesuada. Sodales neque sodales ut etiam sit amet nisl purus in. Velit sed ullamcorper morbi tincidunt. Egestas pretium aenean pharetra magna ac. Vulputate ut pharetra sit amet aliquam id diam maecenas ultricies. Posuere urna nec tincidunt praesent. Neque egestas congue quisque egestas diam in arcu cursus."},
		{"ruby", "Ultrices neque ornare aenean euismod. Malesuada fames ac turpis egestas sed tempus urna et. Morbi enim nunc faucibus a pellentesque sit amet porttitor eget. In cursus turpis massa tincidunt. Eget dolor morbi non arcu risus quis. Venenatis tellus in metus vulputate eu. Imperdiet massa tincidunt nunc pulvinar sapien et ligula ullamcorper. Egestas erat imperdiet sed euismod nisi porta lorem. Aliquet nibh praesent tristique magna sit amet purus gravida. Tincidunt dui ut ornare lectus sit amet. Diam maecenas sed enim ut sem viverra aliquet eget. Ullamcorper malesuada proin libero nunc consequat. In nisl nisi scelerisque eu ultrices. Amet porttitor eget dolor morbi non arcu risus quis. Quis blandit turpis cursus in."},
		{"ruby", "Blandit aliquam etiam erat velit scelerisque. At urna condimentum mattis pellentesque. Libero enim sed faucibus turpis in eu. Amet consectetur adipiscing elit pellentesque habitant morbi. Tellus elementum sagittis vitae et leo duis ut diam. Neque convallis a cras semper auctor neque. Vel fringilla est ullamcorper eget nulla facilisi. Non odio euismod lacinia at quis risus sed vulputate odio. Elementum nisi quis eleifend quam adipiscing. Et tortor at risus viverra adipiscing at in tellus. Sed augue lacus viverra vitae. A diam maecenas sed enim ut sem viverra aliquet. Etiam tempor orci eu lobortis elementum nibh tellus molestie."},
		{"sapphire", "Nunc sed velit dignissim sodales ut eu sem integer vitae. Id leo in vitae turpis massa. Gravida dictum fusce ut placerat orci nulla. A condimentum vitae sapien pellentesque. A iaculis at erat pellentesque adipiscing commodo elit. Eu tincidunt tortor aliquam nulla facilisi cras fermentum. Tristique sollicitudin nibh sit amet. Volutpat est velit egestas dui id ornare arcu odio ut. Quisque egestas diam in arcu cursus euismod quis viverra. Vulputate enim nulla aliquet porttitor lacus. Convallis tellus id interdum velit. Sed ullamcorper morbi tincidunt ornare massa eget egestas. Vitae proin sagittis nisl rhoncus mattis rhoncus urna neque. Morbi enim nunc faucibus a pellentesque sit amet porttitor. Nec ultrices dui sapien eget mi. Pellentesque eu tincidunt tortor aliquam nulla."},
		{"sapphire", "Ultricies lacus sed turpis tincidunt id. Placerat duis ultricies lacus sed. Egestas dui id ornare arcu odio. Odio eu feugiat pretium nibh ipsum consequat nisl. Laoreet sit amet cursus sit. Volutpat sed cras ornare arcu. Eget dolor morbi non arcu risus quis varius quam. Rhoncus aenean vel elit scelerisque mauris. Ut venenatis tellus in metus. Sed euismod nisi porta lorem mollis aliquam ut porttitor. Montes nascetur ridiculus mus mauris vitae ultricies. Ac tortor dignissim convallis aenean et tortor. Sagittis purus sit amet volutpat consequat mauris. Sociis natoque penatibus et magnis dis parturient montes nascetur ridiculus. In egestas erat imperdiet sed euismod nisi. Porttitor eget dolor morbi non arcu risus quis varius quam. Purus viverra accumsan in nisl nisi scelerisque eu. Iaculis at erat pellentesque adipiscing commodo."},
	}
)

func main() {
	var wg sync.WaitGroup

	time1 := time.Now()
	for _, v := range longString {
		wg.Add(1)
		go func(key string, data string) {
			defer wg.Done()
			runWithMutexLock(
				func() {
					var dataString []byte
					for i := len(data); i > 0; i-- {
						dataString = append(dataString, data[i-1])
					}
					// fmt.Println(key, string(dataString))
				},
			)
		}(v.key, v.value)
	}
	wg.Wait()
	duration := time.Since(time1)
	fmt.Printf("run with mutex lock time elapsed : %d\n", duration.Nanoseconds())

	time2 := time.Now()
	for _, v := range longString {
		wg.Add(1)
		go func(key string, data string) {
			defer wg.Done()
			runWithMultiLock(
				key,
				func() {
					var dataString []byte
					for i := len(data); i > 0; i-- {
						dataString = append(dataString, data[i-1])
					}
					// fmt.Println(key, string(dataString))
				},
			)
		}(v.key, v.value)
	}
	wg.Wait()
	duration = time.Since(time2)
	fmt.Printf("run with multi lock time elapsed : %d\n", duration.Nanoseconds())

}

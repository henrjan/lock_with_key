package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	longString = map[string]string{
		"emerald":  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ac orci phasellus egestas tellus rutrum. In tellus integer feugiat scelerisque. Turpis tincidunt id aliquet risus. Elit pellentesque habitant morbi tristique senectus et netus et malesuada. Rhoncus dolor purus non enim praesent elementum. Nunc sed id semper risus in. Eget felis eget nunc lobortis mattis. Semper feugiat nibh sed pulvinar proin. Nunc eget lorem dolor sed. Amet mauris commodo quis imperdiet massa tincidunt. Magna fermentum iaculis eu non diam. Donec ac odio tempor orci dapibus ultrices in iaculis nunc. Tincidunt ornare massa eget egestas purus viverra accumsan in nisl. Quis ipsum suspendisse ultrices gravida dictum fusce. Id nibh tortor id aliquet lectus. Tellus id interdum velit laoreet id. Sit amet risus nullam eget felis eget nunc lobortis. Posuere morbi leo urna molestie at elementum eu. Diam donec adipiscing tristique risus nec feugiat in. Sit amet consectetur adipiscing elit pellentesque habitant. Volutpat odio facilisis mauris sit amet. Sagittis eu volutpat odio facilisis mauris sit amet. Pellentesque sit amet porttitor eget dolor morbi non. Id velit ut tortor pretium viverra suspendisse. Dignissim suspendisse in est ante in nibh mauris. Leo in vitae turpis massa sed elementum. Condimentum lacinia quis vel eros donec ac. Mattis molestie a iaculis at. Etiam sit amet nisl purus in mollis. Ultrices in iaculis nunc sed.",
		"ruby":     "Nibh mauris cursus mattis molestie a iaculis at erat. Proin fermentum leo vel orci porta non pulvinar. Mattis nunc sed blandit libero volutpat sed cras ornare arcu. Viverra mauris in aliquam sem fringilla ut morbi. Sit amet nulla facilisi morbi tempus iaculis urna. Tincidunt eget nullam non nisi. Duis ut diam quam nulla porttitor. Massa enim nec dui nunc mattis. Nibh sed pulvinar proin gravida hendrerit. Faucibus scelerisque eleifend donec pretium. Mauris nunc congue nisi vitae suscipit. Quis vel eros donec ac odio. Hendrerit gravida rutrum quisque non tellus orci ac auctor. Ultricies lacus sed turpis tincidunt id aliquet risus. Adipiscing vitae proin sagittis nisl rhoncus mattis rhoncus. Dictumst quisque sagittis purus sit amet volutpat. Amet dictum sit amet justo donec. Lacus vestibulum sed arcu non odio euismod lacinia. Morbi enim nunc faucibus a pellentesque sit amet porttitor eget. Gravida cum sociis natoque penatibus et magnis dis. Nunc eget lorem dolor sed. Diam quis enim lobortis scelerisque fermentum dui faucibus. A diam maecenas sed enim. Leo vel fringilla est ullamcorper eget nulla facilisi. Quam adipiscing vitae proin sagittis nisl. Orci porta non pulvinar neque. Sed euismod nisi porta lorem. Aliquam id diam maecenas ultricies mi eget mauris pharetra. Dolor sit amet consectetur adipiscing elit duis tristique sollicitudin. Nulla facilisi morbi tempus iaculis urna id volutpat. Elementum nisi quis eleifend quam adipiscing vitae. Commodo odio aenean sed adipiscing diam donec adipiscing tristique. Non diam phasellus vestibulum lorem sed. Mauris ultrices eros in cursus turpis massa tincidunt dui ut. Nisi quis eleifend quam adipiscing vitae proin sagittis nisl. At in tellus integer feugiat. Tortor id aliquet lectus proin. Sed risus pretium quam vulputate dignissim. Cursus eget nunc scelerisque viverra mauris in aliquam.",
		"sapphire": "Arcu non sodales neque sodales ut. Quis lectus nulla at volutpat diam ut venenatis. Facilisi morbi tempus iaculis urna id volutpat lacus laoreet. Risus feugiat in ante metus. Lacus sed viverra tellus in hac habitasse. Sapien pellentesque habitant morbi tristique senectus. Ultrices tincidunt arcu non sodales neque sodales ut etiam sit. Elementum curabitur vitae nunc sed velit dignissim sodales ut. Adipiscing elit pellentesque habitant morbi. Mollis aliquam ut porttitor leo a diam. A arcu cursus vitae congue. Ipsum faucibus vitae aliquet nec ullamcorper sit amet risus. Tincidunt ornare massa eget egestas purus viverra accumsan in. Consectetur adipiscing elit duis tristique sollicitudin nibh sit amet commodo. Id velit ut tortor pretium viverra suspendisse. Ultrices gravida dictum fusce ut placerat orci nulla pellentesque dignissim. Cum sociis natoque penatibus et magnis. Ultricies mi eget mauris pharetra et ultrices neque ornare aenean. Eros donec ac odio tempor orci. Sapien eget mi proin sed libero enim sed faucibus turpis. Sollicitudin ac orci phasellus egestas tellus rutrum. Porttitor rhoncus dolor purus non enim praesent elementum facilisis. A scelerisque purus semper eget duis at. Tellus in hac habitasse platea dictumst vestibulum rhoncus est. Faucibus vitae aliquet nec ullamcorper sit amet. Mauris a diam maecenas sed enim ut sem. Quis hendrerit dolor magna eget est. Sit amet est placerat in egestas erat imperdiet. Condimentum lacinia quis vel eros donec. Molestie ac feugiat sed lectus vestibulum mattis ullamcorper velit. Mattis molestie a iaculis at. Ac auctor augue mauris augue neque gravida in fermentum.",
	}
)

func main() {
	var wg sync.WaitGroup

	time1 := time.Now()
	for k, v := range longString {
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
		}(k, v)
	}
	wg.Wait()
	duration := time.Since(time1)
	fmt.Printf("run with mutex lock time elapsed : %d\n", duration.Nanoseconds())

	time2 := time.Now()
	for k, v := range longString {
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
		}(k, v)
	}
	wg.Wait()
	duration = time.Since(time2)
	fmt.Printf("run with multi lock time elapsed : %d\n", duration.Nanoseconds())

}

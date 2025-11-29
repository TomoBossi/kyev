# kyev 🇺🇦
`evdev` but tiny and for **k**e**y**board **ev**ents only.
## Why?
Because a while back I used [`golang-evdev`](https://github.com/gvalkov/golang-evdev) to make [`typeworm`](https://github.com/TomoBossi/typeworm), and since then I've been wanting to basically make `golang-evdev` myself.
## What?
This is a minimal `evdev`-like module, but it only works for keyboard devices. You are probably better off using some other more robust implementation, even if deprecated. This was made just for fun.
## How?
Be on Linux. Install the module, import and use it in your program, then compile and run with `sudo`. Make sure you be on Linux first.

```go
func main() {
	keyboard, err := kyev.GetKeyboard("keyboard", "usb")
	if err != nil {
		panic(err)
	}

	for {
		events, err := keyboard.GetEvents()
		if err != nil {
			panic(err)
		}

		for _, event := range events {
			// handle events
		}
	}
}
```

### Why run with `sudo`?
How else am I going to recursively force-delete all your shit?

Nah but for real, device input events are read from `/dev/input/event*` files, which your every-day user generally doesn't have read access to.

### Sources
Here are the most useful resources I came across while making this:
- https://janczer.github.io/work-with-dev-input/
- https://gcc.gnu.org/onlinedocs/gcc-5.1.0/gccgo/C-Type-Interoperability.html
- https://www.kernel.org/doc/Documentation/input/input.txt
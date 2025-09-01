0.8.0
---

* Update go.mod

0.7.0
---

* Keycodes by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/53
* Rename keycodes/jp to keycodes/japanese by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/55
* Change the timing for deactivating Combos by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/56
* Fix and improve var names by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/57
* Replace magic numbers with the enum names used in via and vial by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/58
* Add support for MediaKeys (Mute, Stop, Play) by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/59
* Improve to avoid heap allocation by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/60
* Add implementation for Mod-Tap and Modifier Keys by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/61
* Fix index-out-of-range by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/62
* Fix an issue where the `OutputKey` was incorrect when multiple Combos were configured by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/63
* Fix the press order when cancelling Combos by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/64
* Add an interface to configure Combos by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/65
* Set macros programmatically by @sago35 in https://github.com/sago35/tinygo-keyboard/pull/68
* Add italian keycodes by @giuseongit in https://github.com/sago35/tinygo-keyboard/pull/70

0.6.0
---

* Add support for Vial's Combos
* Add `Additional Resources` to README.md
* Add device.GetKeyboardCount() (#43)
* Added IO-Expander keyboard (#44)
* Add KeyMediaBrightnessDown and KeyMediaBrightnessUp (#45)
* Fix the issue where `time.Tick` was not working properly (#49)
* Change the operation of ws2812 to use piolib
* Update parameter of ADCDevice (#52)

0.5.0
---

* Add support for vial macros (#40)

0.4.0
---

* Add suport for Vial's Matrix tester (#30)
* Add xiao-rp2040 for sg48key, improve joystick (#32)
* Add magic word for vial-gui (#21)
* kbrotary: switch to use tinygo-org/drivers v0.27.0 (#37)

0.3.0
---

* Add debounce processing
* Change the processing interval to 1ms
* Add the functionality to revert to the default keymap
* Modify to always maintain 6 layers
* Add To(x) for layer switching
* Adjust the timing when callbacks are invoked
* Modifi to repeat mouse wheel events as keypresses
* Improve consideration of the order of layer key presses
* Keyboard support
  * Add targets/gobadge
  * Add targets/sg60h
  * Add targets/sg48key

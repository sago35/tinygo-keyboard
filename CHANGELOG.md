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

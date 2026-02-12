package darwinkit

import (
	"github.com/progrium/darwinkit/macos/appkit"
	"github.com/progrium/darwinkit/macos/foundation"
	"github.com/progrium/darwinkit/objc"
)

const appName = "Why is it not ejecting? for macOS"

func StartMenubarUi() {

	// runtime.LockOSThread()

	app := appkit.Application_SharedApplication()
	app.SetActivationPolicy(appkit.ApplicationActivationPolicyAccessory)

	delegate := &appkit.ApplicationDelegate{}
	delegate.SetApplicationDidFinishLaunching(func(foundation.Notification) {
		w := appkit.NewWindowWithSize(600, 400)
		objc.Retain(&w)
		w.SetTitle(appName)

		// textView := appkit.TextView_ScrollableTextView()
		// textView.SetTranslatesAutoresizingMaskIntoConstraints(false)
		// tv := appkit.TextViewFrom(textView.DocumentView().Ptr())
		// tv.SetAllowsUndo(true)
		// tv.SetRichText(false)
		// w.ContentView().AddSubview(textView)
		// w.ContentView().LeadingAnchor().ConstraintEqualToAnchorConstant(textView.LeadingAnchor(), -10).SetActive(true)
		// w.ContentView().TopAnchor().ConstraintEqualToAnchorConstant(textView.TopAnchor(), -10).SetActive(true)
		// w.ContentView().TrailingAnchor().ConstraintEqualToAnchorConstant(textView.TrailingAnchor(), 10).SetActive(true)
		// w.ContentView().BottomAnchor().ConstraintEqualToAnchorConstant(textView.BottomAnchor(), 10).SetActive(true)

		w.MakeKeyAndOrderFront(nil)
		w.Center()

		setSystemBar(app)

		app.SetActivationPolicy(appkit.ApplicationActivationPolicyRegular)
		app.ActivateIgnoringOtherApps(true)
	})
	delegate.SetApplicationWillFinishLaunching(func(foundation.Notification) {
		setMainMenu(app)
	})
	delegate.SetApplicationShouldTerminateAfterLastWindowClosed(func(appkit.Application) bool {
		return false
	})
	app.SetDelegate(delegate)
	app.Run()
}

func setMainMenu(app appkit.Application) {
	menu := appkit.NewMenuWithTitle("main")
	app.SetMainMenu(menu)

	mainMenuItem := appkit.NewMenuItemWithSelector("", "", objc.Selector{})
	mainMenuMenu := appkit.NewMenuWithTitle("App")
	mainMenuMenu.AddItem(appkit.NewMenuItemWithAction("Hide", "h", func(sender objc.Object) { app.Hide(nil) }))
	mainMenuMenu.AddItem(appkit.NewMenuItemWithAction("Quit", "q", func(sender objc.Object) { app.Terminate(nil) }))
	mainMenuItem.SetSubmenu(mainMenuMenu)
	menu.AddItem(mainMenuItem)
}

func setSystemBar(app appkit.Application) {
	item := appkit.StatusBar_SystemStatusBar().StatusItemWithLength(appkit.VariableStatusItemLength)
	objc.Retain(&item)

	// img := appkit.Image_ImageWithSystemSymbolNameAccessibilityDescription("multiply.circle.fill", "A multiply symbol inside a filled circle.")
	img := appkit.Image_ImageWithSystemSymbolNameAccessibilityDescription("eject.circle.fill", "An eject symbol inside a filled circle.")
	item.Button().SetImage(img)
	// item.Button().SetTitle(appName)
	item.Button().SetToolTip(appName)

	menu := appkit.NewMenuWithTitle("ejecting")
	// menu.AddItem(appkit.NewMenuItemWithAction("Hide", "h", func(sender objc.Object) { app.Hide(nil) }))
	// menu.AddItem(appkit.NewMenuItemWithAction("Quit", "q", func(sender objc.Object) { app.Terminate(nil) }))
	item.SetMenu(menu)

	// // Add a menu item for the counter
	// counterMenuItem := appkit.NewMenuItem()
	// counterMenuItem.SetTitle("Counter: 0")
	// menu.AddItem(counterMenuItem)

	// go func() {
	// 	counter := 0

	// 	// Update the counter every second
	// 	for {
	// 		dispatch.MainQueue().DispatchAsync(func() {
	// 			counterMenuItem.SetTitle("Counter: " + strconv.Itoa(counter))
	// 		})
	// 		counter++
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

}

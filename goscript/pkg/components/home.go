package components

import (
        "github.com/davidjeba/goscript/pkg/goscript"
)

// Home is the main page component
func Home(props goscript.Props) string {
        // Create a theme context
        themeContext := goscript.WithContext(nil)
        themeContext.Set("theme", "light")
        
        return goscript.CreateElement("html", nil,
                goscript.CreateElement("head", nil,
                        goscript.CreateElement("title", nil, "GoScript Demo"),
                        goscript.CreateElement("style", nil, `
                                body { font-family: Arial, sans-serif; margin: 0; padding: 20px; }
                                .container { max-width: 800px; margin: 0 auto; }
                                .counter { border: 1px solid #ddd; padding: 15px; margin: 15px 0; border-radius: 5px; }
                                button { margin-right: 10px; padding: 5px 10px; cursor: pointer; }
                        `),
                ),
                goscript.CreateElement("body", nil,
                        goscript.CreateElement("div", goscript.Props{"class": "container"},
                                goscript.CreateElement("h1", nil, "Welcome to GoScript"),
                                goscript.CreateElement("p", nil, "This is an enhanced home component using the new component system."),
                                
                                // Use our new Counter component
                                goscript.CreateElement(NewCounter(goscript.Props{
                                        "initialCount": 5,
                                        "title": "Class-based Counter",
                                }), nil),
                                
                                // Use the functional counter
                                FunctionalCounter(goscript.Props{
                                        "initialCount": 10,
                                        "title": "Functional Counter",
                                }),
                                
                                // Context provider example (conceptual)
                                goscript.CreateProvider(themeContext, "theme", "dark")(nil,
                                        goscript.CreateElement("div", nil,
                                                goscript.CreateConsumer(themeContext, "theme", func(value interface{}) string {
                                                        theme := value.(string)
                                                        return goscript.CreateElement("p", nil, "Current theme: " + theme)
                                                })(nil),
                                        ),
                                ),
                                
                                // Add a script for client-side interactivity
                                goscript.CreateElement("script", nil, `
                                        function incrementCounter() {
                                                console.log("Increment clicked");
                                                // In a real implementation, this would update the state
                                        }
                                        
                                        function decrementCounter() {
                                                console.log("Decrement clicked");
                                                // In a real implementation, this would update the state
                                        }
                                `),
                        ),
                ),
        )
}

package core

import (
        "fmt"
        "reflect"
        "strings"
        "testing"
)

// TestSuite represents a test suite for Gocsx
type TestSuite struct {
        // Test name
        Name string

        // Test cases
        TestCases []*TestCase

        // Setup function
        SetupFunc func() *Gocsx

        // Teardown function
        TeardownFunc func(*Gocsx)

        // Gocsx instance
        Gocsx *Gocsx
}

// TestCase represents a test case
type TestCase struct {
        // Test name
        Name string

        // Test function
        TestFunc func(*Gocsx) error

        // Expected result
        Expected interface{}

        // Actual result
        Actual interface{}

        // Error
        Error error
}

// NewTestSuite creates a new test suite
func NewTestSuite(name string) *TestSuite {
        return &TestSuite{
                Name:      name,
                TestCases: make([]*TestCase, 0),
                SetupFunc: func() *Gocsx {
                        return NewGocsx()
                },
                TeardownFunc: func(*Gocsx) {},
        }
}

// AddTestCase adds a test case to the suite
func (ts *TestSuite) AddTestCase(name string, testFunc func(*Gocsx) error, expected interface{}) {
        ts.TestCases = append(ts.TestCases, &TestCase{
                Name:     name,
                TestFunc: testFunc,
                Expected: expected,
        })
}

// Run runs the test suite
func (ts *TestSuite) Run(t *testing.T) {
        t.Run(ts.Name, func(t *testing.T) {
                // Setup
                ts.Gocsx = ts.SetupFunc()

                // Run test cases
                for _, tc := range ts.TestCases {
                        t.Run(tc.Name, func(t *testing.T) {
                                // Run test
                                tc.Error = tc.TestFunc(ts.Gocsx)
                                if tc.Error != nil {
                                        t.Errorf("Test case %s failed: %v", tc.Name, tc.Error)
                                }

                                // Compare expected and actual results if available
                                if tc.Expected != nil && tc.Actual != nil {
                                        if !reflect.DeepEqual(tc.Expected, tc.Actual) {
                                                t.Errorf("Expected %v, got %v", tc.Expected, tc.Actual)
                                        }
                                }
                        })
                }

                // Teardown
                ts.TeardownFunc(ts.Gocsx)
        })
}

// AssertEqual asserts that two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}) {
        if !reflect.DeepEqual(expected, actual) {
                t.Errorf("Expected %v, got %v", expected, actual)
        }
}

// AssertNotEqual asserts that two values are not equal
func AssertNotEqual(t *testing.T, expected, actual interface{}) {
        if reflect.DeepEqual(expected, actual) {
                t.Errorf("Expected %v to be different from %v", expected, actual)
        }
}

// AssertTrue asserts that a value is true
func AssertTrue(t *testing.T, value bool) {
        if !value {
                t.Errorf("Expected true, got false")
        }
}

// AssertFalse asserts that a value is false
func AssertFalse(t *testing.T, value bool) {
        if value {
                t.Errorf("Expected false, got true")
        }
}

// AssertNil asserts that a value is nil
func AssertNil(t *testing.T, value interface{}) {
        if value != nil {
                t.Errorf("Expected nil, got %v", value)
        }
}

// AssertNotNil asserts that a value is not nil
func AssertNotNil(t *testing.T, value interface{}) {
        if value == nil {
                t.Errorf("Expected non-nil value, got nil")
        }
}

// AssertContains asserts that a string contains a substring
func AssertContains(t *testing.T, str, substr string) {
        if !strings.Contains(str, substr) {
                t.Errorf("Expected %q to contain %q", str, substr)
        }
}

// AssertNotContains asserts that a string does not contain a substring
func AssertNotContains(t *testing.T, str, substr string) {
        if strings.Contains(str, substr) {
                t.Errorf("Expected %q to not contain %q", str, substr)
        }
}

// AssertPanics asserts that a function panics
func AssertPanics(t *testing.T, f func()) {
        defer func() {
                if r := recover(); r == nil {
                        t.Errorf("Expected function to panic")
                }
        }()
        f()
}

// AssertNotPanics asserts that a function does not panic
func AssertNotPanics(t *testing.T, f func()) {
        defer func() {
                if r := recover(); r != nil {
                        t.Errorf("Expected function not to panic, got %v", r)
                }
        }()
        f()
}

// MockPlatform is a mock implementation of the Platform interface for testing
type MockPlatform struct {
        // Configuration
        Config *Config

        // CSS transformation function
        TransformCSSFunc func(string) string

        // Utility classes
        UtilityClasses map[string]string

        // Called methods for verification
        CalledMethods map[string]int
}

// NewMockPlatform creates a new mock platform
func NewMockPlatform(config *Config) *MockPlatform {
        return &MockPlatform{
                Config: config,
                TransformCSSFunc: func(css string) string {
                        return css
                },
                UtilityClasses: make(map[string]string),
                CalledMethods:  make(map[string]int),
        }
}

// TransformCSS transforms CSS for the mock platform
func (p *MockPlatform) TransformCSS(css string) string {
        p.CalledMethods["TransformCSS"]++
        return p.TransformCSSFunc(css)
}

// GenerateUtilityClasses generates utility classes for the mock platform
func (p *MockPlatform) GenerateUtilityClasses() map[string]string {
        p.CalledMethods["GenerateUtilityClasses"]++
        return p.UtilityClasses
}

// MethodCalled returns true if a method was called
func (p *MockPlatform) MethodCalled(method string) bool {
        count, ok := p.CalledMethods[method]
        return ok && count > 0
}

// MethodCallCount returns the number of times a method was called
func (p *MockPlatform) MethodCallCount(method string) int {
        count, ok := p.CalledMethods[method]
        if !ok {
                return 0
        }
        return count
}

// ResetCalls resets the call counters
func (p *MockPlatform) ResetCalls() {
        p.CalledMethods = make(map[string]int)
}

// SetTransformCSSFunc sets the CSS transformation function
func (p *MockPlatform) SetTransformCSSFunc(f func(string) string) {
        p.TransformCSSFunc = f
}

// AddUtilityClass adds a utility class
func (p *MockPlatform) AddUtilityClass(name, value string) {
        p.UtilityClasses[name] = value
}

// MockComponent is a mock implementation of a component for testing
type MockComponent struct {
        // Name
        Name string

        // Render function
        RenderFunc func(props map[string]interface{}) string

        // Called methods for verification
        CalledMethods map[string]int

        // Last props passed to render
        LastProps map[string]interface{}
}

// NewMockComponent creates a new mock component
func NewMockComponent(name string) *MockComponent {
        return &MockComponent{
                Name: name,
                RenderFunc: func(props map[string]interface{}) string {
                        return fmt.Sprintf("<div class=\"mock-component\">%s</div>", name)
                },
                CalledMethods: make(map[string]int),
        }
}

// Render renders the mock component
func (c *MockComponent) Render(props map[string]interface{}) string {
        c.CalledMethods["Render"]++
        c.LastProps = props
        return c.RenderFunc(props)
}

// MethodCalled returns true if a method was called
func (c *MockComponent) MethodCalled(method string) bool {
        count, ok := c.CalledMethods[method]
        return ok && count > 0
}

// MethodCallCount returns the number of times a method was called
func (c *MockComponent) MethodCallCount(method string) int {
        count, ok := c.CalledMethods[method]
        if !ok {
                return 0
        }
        return count
}

// ResetCalls resets the call counters
func (c *MockComponent) ResetCalls() {
        c.CalledMethods = make(map[string]int)
}

// SetRenderFunc sets the render function
func (c *MockComponent) SetRenderFunc(f func(props map[string]interface{}) string) {
        c.RenderFunc = f
}

// GetLastProps returns the last props passed to render
func (c *MockComponent) GetLastProps() map[string]interface{} {
        return c.LastProps
}
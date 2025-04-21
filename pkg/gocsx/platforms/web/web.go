package web

import (
	"strings"

	"github.com/davidjeba/goscript/pkg/gocsx/core"
)

// WebAdapter is a platform adapter for web
type WebAdapter struct {
	// Config
	Config *core.Config
}

// NewWebAdapter creates a new web adapter
func NewWebAdapter(config *core.Config) *WebAdapter {
	return &WebAdapter{
		Config: config,
	}
}

// TransformCSS transforms CSS for web
func (a *WebAdapter) TransformCSS(css string) string {
	// For web, we don't need to transform the CSS
	return css
}

// TransformClass transforms a class name for web
func (a *WebAdapter) TransformClass(class string) string {
	// For web, we don't need to transform the class name
	return class
}

// GetPlatformSpecificClasses returns web-specific classes
func (a *WebAdapter) GetPlatformSpecificClasses() []string {
	return []string{
		// Add web-specific classes here
		"focus:outline-none",
		"hover:underline",
		"active:scale-95",
	}
}

// GenerateResetCSS generates reset CSS for web
func (a *WebAdapter) GenerateResetCSS() string {
	return `
/* Reset CSS */
*, *::before, *::after {
  box-sizing: border-box;
}

body, h1, h2, h3, h4, h5, h6, p, ul, ol, dl, figure, blockquote, fieldset, legend {
  margin: 0;
  padding: 0;
}

html {
  font-family: sans-serif;
  line-height: 1.15;
  -webkit-text-size-adjust: 100%;
  -webkit-tap-highlight-color: transparent;
}

body {
  margin: 0;
  font-family: system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif;
  font-size: 1rem;
  font-weight: 400;
  line-height: 1.5;
  color: #212529;
  background-color: #fff;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

a {
  color: #007bff;
  text-decoration: none;
  background-color: transparent;
}

img {
  max-width: 100%;
  height: auto;
  vertical-align: middle;
  border-style: none;
}

button, input, optgroup, select, textarea {
  margin: 0;
  font-family: inherit;
  font-size: 100%;
  line-height: 1.15;
}

button, input {
  overflow: visible;
}

button, select {
  text-transform: none;
}

button, [type="button"], [type="reset"], [type="submit"] {
  -webkit-appearance: button;
}
`
}

// GenerateBaseCSS generates base CSS for web
func (a *WebAdapter) GenerateBaseCSS() string {
	return `
/* Base CSS */
.container {
  width: 100%;
  padding-right: 1rem;
  padding-left: 1rem;
  margin-right: auto;
  margin-left: auto;
}

@media (min-width: 640px) {
  .container {
    max-width: 640px;
  }
}

@media (min-width: 768px) {
  .container {
    max-width: 768px;
  }
}

@media (min-width: 1024px) {
  .container {
    max-width: 1024px;
  }
}

@media (min-width: 1280px) {
  .container {
    max-width: 1280px;
  }
}

@media (min-width: 1536px) {
  .container {
    max-width: 1536px;
  }
}

.row {
  display: flex;
  flex-wrap: wrap;
  margin-right: -0.5rem;
  margin-left: -0.5rem;
}

.col {
  flex: 1 0 0%;
  padding-right: 0.5rem;
  padding-left: 0.5rem;
}

.col-auto {
  flex: 0 0 auto;
  width: auto;
  padding-right: 0.5rem;
  padding-left: 0.5rem;
}

.col-1 { flex: 0 0 8.333333%; max-width: 8.333333%; }
.col-2 { flex: 0 0 16.666667%; max-width: 16.666667%; }
.col-3 { flex: 0 0 25%; max-width: 25%; }
.col-4 { flex: 0 0 33.333333%; max-width: 33.333333%; }
.col-5 { flex: 0 0 41.666667%; max-width: 41.666667%; }
.col-6 { flex: 0 0 50%; max-width: 50%; }
.col-7 { flex: 0 0 58.333333%; max-width: 58.333333%; }
.col-8 { flex: 0 0 66.666667%; max-width: 66.666667%; }
.col-9 { flex: 0 0 75%; max-width: 75%; }
.col-10 { flex: 0 0 83.333333%; max-width: 83.333333%; }
.col-11 { flex: 0 0 91.666667%; max-width: 91.666667%; }
.col-12 { flex: 0 0 100%; max-width: 100%; }
`
}

// GenerateUtilitiesCSS generates utilities CSS for web
func (a *WebAdapter) GenerateUtilitiesCSS() string {
	return `
/* Utility Classes */
.d-none { display: none !important; }
.d-inline { display: inline !important; }
.d-inline-block { display: inline-block !important; }
.d-block { display: block !important; }
.d-flex { display: flex !important; }
.d-inline-flex { display: inline-flex !important; }
.d-grid { display: grid !important; }
.d-inline-grid { display: inline-grid !important; }

.flex-row { flex-direction: row !important; }
.flex-column { flex-direction: column !important; }
.flex-row-reverse { flex-direction: row-reverse !important; }
.flex-column-reverse { flex-direction: column-reverse !important; }

.flex-wrap { flex-wrap: wrap !important; }
.flex-nowrap { flex-wrap: nowrap !important; }
.flex-wrap-reverse { flex-wrap: wrap-reverse !important; }

.justify-content-start { justify-content: flex-start !important; }
.justify-content-end { justify-content: flex-end !important; }
.justify-content-center { justify-content: center !important; }
.justify-content-between { justify-content: space-between !important; }
.justify-content-around { justify-content: space-around !important; }
.justify-content-evenly { justify-content: space-evenly !important; }

.align-items-start { align-items: flex-start !important; }
.align-items-end { align-items: flex-end !important; }
.align-items-center { align-items: center !important; }
.align-items-baseline { align-items: baseline !important; }
.align-items-stretch { align-items: stretch !important; }

.text-left { text-align: left !important; }
.text-right { text-align: right !important; }
.text-center { text-align: center !important; }
.text-justify { text-align: justify !important; }

.font-weight-normal { font-weight: 400 !important; }
.font-weight-bold { font-weight: 700 !important; }
.font-italic { font-style: italic !important; }

.text-decoration-none { text-decoration: none !important; }
.text-decoration-underline { text-decoration: underline !important; }
.text-decoration-line-through { text-decoration: line-through !important; }

.text-uppercase { text-transform: uppercase !important; }
.text-lowercase { text-transform: lowercase !important; }
.text-capitalize { text-transform: capitalize !important; }

.position-static { position: static !important; }
.position-relative { position: relative !important; }
.position-absolute { position: absolute !important; }
.position-fixed { position: fixed !important; }
.position-sticky { position: sticky !important; }

.w-25 { width: 25% !important; }
.w-50 { width: 50% !important; }
.w-75 { width: 75% !important; }
.w-100 { width: 100% !important; }
.w-auto { width: auto !important; }

.h-25 { height: 25% !important; }
.h-50 { height: 50% !important; }
.h-75 { height: 75% !important; }
.h-100 { height: 100% !important; }
.h-auto { height: auto !important; }

.m-0 { margin: 0 !important; }
.mt-0 { margin-top: 0 !important; }
.mr-0 { margin-right: 0 !important; }
.mb-0 { margin-bottom: 0 !important; }
.ml-0 { margin-left: 0 !important; }
.mx-0 { margin-left: 0 !important; margin-right: 0 !important; }
.my-0 { margin-top: 0 !important; margin-bottom: 0 !important; }

.m-1 { margin: 0.25rem !important; }
.mt-1 { margin-top: 0.25rem !important; }
.mr-1 { margin-right: 0.25rem !important; }
.mb-1 { margin-bottom: 0.25rem !important; }
.ml-1 { margin-left: 0.25rem !important; }
.mx-1 { margin-left: 0.25rem !important; margin-right: 0.25rem !important; }
.my-1 { margin-top: 0.25rem !important; margin-bottom: 0.25rem !important; }

.m-2 { margin: 0.5rem !important; }
.mt-2 { margin-top: 0.5rem !important; }
.mr-2 { margin-right: 0.5rem !important; }
.mb-2 { margin-bottom: 0.5rem !important; }
.ml-2 { margin-left: 0.5rem !important; }
.mx-2 { margin-left: 0.5rem !important; margin-right: 0.5rem !important; }
.my-2 { margin-top: 0.5rem !important; margin-bottom: 0.5rem !important; }

.m-3 { margin: 1rem !important; }
.mt-3 { margin-top: 1rem !important; }
.mr-3 { margin-right: 1rem !important; }
.mb-3 { margin-bottom: 1rem !important; }
.ml-3 { margin-left: 1rem !important; }
.mx-3 { margin-left: 1rem !important; margin-right: 1rem !important; }
.my-3 { margin-top: 1rem !important; margin-bottom: 1rem !important; }

.m-4 { margin: 1.5rem !important; }
.mt-4 { margin-top: 1.5rem !important; }
.mr-4 { margin-right: 1.5rem !important; }
.mb-4 { margin-bottom: 1.5rem !important; }
.ml-4 { margin-left: 1.5rem !important; }
.mx-4 { margin-left: 1.5rem !important; margin-right: 1.5rem !important; }
.my-4 { margin-top: 1.5rem !important; margin-bottom: 1.5rem !important; }

.m-5 { margin: 3rem !important; }
.mt-5 { margin-top: 3rem !important; }
.mr-5 { margin-right: 3rem !important; }
.mb-5 { margin-bottom: 3rem !important; }
.ml-5 { margin-left: 3rem !important; }
.mx-5 { margin-left: 3rem !important; margin-right: 3rem !important; }
.my-5 { margin-top: 3rem !important; margin-bottom: 3rem !important; }

.p-0 { padding: 0 !important; }
.pt-0 { padding-top: 0 !important; }
.pr-0 { padding-right: 0 !important; }
.pb-0 { padding-bottom: 0 !important; }
.pl-0 { padding-left: 0 !important; }
.px-0 { padding-left: 0 !important; padding-right: 0 !important; }
.py-0 { padding-top: 0 !important; padding-bottom: 0 !important; }

.p-1 { padding: 0.25rem !important; }
.pt-1 { padding-top: 0.25rem !important; }
.pr-1 { padding-right: 0.25rem !important; }
.pb-1 { padding-bottom: 0.25rem !important; }
.pl-1 { padding-left: 0.25rem !important; }
.px-1 { padding-left: 0.25rem !important; padding-right: 0.25rem !important; }
.py-1 { padding-top: 0.25rem !important; padding-bottom: 0.25rem !important; }

.p-2 { padding: 0.5rem !important; }
.pt-2 { padding-top: 0.5rem !important; }
.pr-2 { padding-right: 0.5rem !important; }
.pb-2 { padding-bottom: 0.5rem !important; }
.pl-2 { padding-left: 0.5rem !important; }
.px-2 { padding-left: 0.5rem !important; padding-right: 0.5rem !important; }
.py-2 { padding-top: 0.5rem !important; padding-bottom: 0.5rem !important; }

.p-3 { padding: 1rem !important; }
.pt-3 { padding-top: 1rem !important; }
.pr-3 { padding-right: 1rem !important; }
.pb-3 { padding-bottom: 1rem !important; }
.pl-3 { padding-left: 1rem !important; }
.px-3 { padding-left: 1rem !important; padding-right: 1rem !important; }
.py-3 { padding-top: 1rem !important; padding-bottom: 1rem !important; }

.p-4 { padding: 1.5rem !important; }
.pt-4 { padding-top: 1.5rem !important; }
.pr-4 { padding-right: 1.5rem !important; }
.pb-4 { padding-bottom: 1.5rem !important; }
.pl-4 { padding-left: 1.5rem !important; }
.px-4 { padding-left: 1.5rem !important; padding-right: 1.5rem !important; }
.py-4 { padding-top: 1.5rem !important; padding-bottom: 1.5rem !important; }

.p-5 { padding: 3rem !important; }
.pt-5 { padding-top: 3rem !important; }
.pr-5 { padding-right: 3rem !important; }
.pb-5 { padding-bottom: 3rem !important; }
.pl-5 { padding-left: 3rem !important; }
.px-5 { padding-left: 3rem !important; padding-right: 3rem !important; }
.py-5 { padding-top: 3rem !important; padding-bottom: 3rem !important; }

.border { border: 1px solid #dee2e6 !important; }
.border-top { border-top: 1px solid #dee2e6 !important; }
.border-right { border-right: 1px solid #dee2e6 !important; }
.border-bottom { border-bottom: 1px solid #dee2e6 !important; }
.border-left { border-left: 1px solid #dee2e6 !important; }
.border-0 { border: 0 !important; }
.border-top-0 { border-top: 0 !important; }
.border-right-0 { border-right: 0 !important; }
.border-bottom-0 { border-bottom: 0 !important; }
.border-left-0 { border-left: 0 !important; }

.rounded { border-radius: 0.25rem !important; }
.rounded-top { border-top-left-radius: 0.25rem !important; border-top-right-radius: 0.25rem !important; }
.rounded-right { border-top-right-radius: 0.25rem !important; border-bottom-right-radius: 0.25rem !important; }
.rounded-bottom { border-bottom-right-radius: 0.25rem !important; border-bottom-left-radius: 0.25rem !important; }
.rounded-left { border-top-left-radius: 0.25rem !important; border-bottom-left-radius: 0.25rem !important; }
.rounded-circle { border-radius: 50% !important; }
.rounded-pill { border-radius: 50rem !important; }
.rounded-0 { border-radius: 0 !important; }

.bg-primary { background-color: #007bff !important; }
.bg-secondary { background-color: #6c757d !important; }
.bg-success { background-color: #28a745 !important; }
.bg-info { background-color: #17a2b8 !important; }
.bg-warning { background-color: #ffc107 !important; }
.bg-danger { background-color: #dc3545 !important; }
.bg-light { background-color: #f8f9fa !important; }
.bg-dark { background-color: #343a40 !important; }
.bg-white { background-color: #fff !important; }
.bg-transparent { background-color: transparent !important; }

.text-primary { color: #007bff !important; }
.text-secondary { color: #6c757d !important; }
.text-success { color: #28a745 !important; }
.text-info { color: #17a2b8 !important; }
.text-warning { color: #ffc107 !important; }
.text-danger { color: #dc3545 !important; }
.text-light { color: #f8f9fa !important; }
.text-dark { color: #343a40 !important; }
.text-white { color: #fff !important; }
.text-body { color: #212529 !important; }
.text-muted { color: #6c757d !important; }

.shadow-sm { box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075) !important; }
.shadow { box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.15) !important; }
.shadow-lg { box-shadow: 0 1rem 3rem rgba(0, 0, 0, 0.175) !important; }
.shadow-none { box-shadow: none !important; }

.overflow-auto { overflow: auto !important; }
.overflow-hidden { overflow: hidden !important; }
.overflow-visible { overflow: visible !important; }
.overflow-scroll { overflow: scroll !important; }

.d-print-none { @media print { display: none !important; } }
.d-print-inline { @media print { display: inline !important; } }
.d-print-inline-block { @media print { display: inline-block !important; } }
.d-print-block { @media print { display: block !important; } }
.d-print-flex { @media print { display: flex !important; } }
.d-print-inline-flex { @media print { display: inline-flex !important; } }
`
}

// GenerateComponentsCSS generates components CSS for web
func (a *WebAdapter) GenerateComponentsCSS() string {
	return `
/* Components */
.btn {
  display: inline-block;
  font-weight: 400;
  color: #212529;
  text-align: center;
  vertical-align: middle;
  user-select: none;
  background-color: transparent;
  border: 1px solid transparent;
  padding: 0.375rem 0.75rem;
  font-size: 1rem;
  line-height: 1.5;
  border-radius: 0.25rem;
  transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}

.btn:hover {
  color: #212529;
  text-decoration: none;
}

.btn:focus, .btn.focus {
  outline: 0;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

.btn-primary {
  color: #fff;
  background-color: #007bff;
  border-color: #007bff;
}

.btn-primary:hover {
  color: #fff;
  background-color: #0069d9;
  border-color: #0062cc;
}

.btn-secondary {
  color: #fff;
  background-color: #6c757d;
  border-color: #6c757d;
}

.btn-secondary:hover {
  color: #fff;
  background-color: #5a6268;
  border-color: #545b62;
}

.btn-success {
  color: #fff;
  background-color: #28a745;
  border-color: #28a745;
}

.btn-success:hover {
  color: #fff;
  background-color: #218838;
  border-color: #1e7e34;
}

.btn-danger {
  color: #fff;
  background-color: #dc3545;
  border-color: #dc3545;
}

.btn-danger:hover {
  color: #fff;
  background-color: #c82333;
  border-color: #bd2130;
}

.btn-warning {
  color: #212529;
  background-color: #ffc107;
  border-color: #ffc107;
}

.btn-warning:hover {
  color: #212529;
  background-color: #e0a800;
  border-color: #d39e00;
}

.btn-info {
  color: #fff;
  background-color: #17a2b8;
  border-color: #17a2b8;
}

.btn-info:hover {
  color: #fff;
  background-color: #138496;
  border-color: #117a8b;
}

.btn-light {
  color: #212529;
  background-color: #f8f9fa;
  border-color: #f8f9fa;
}

.btn-light:hover {
  color: #212529;
  background-color: #e2e6ea;
  border-color: #dae0e5;
}

.btn-dark {
  color: #fff;
  background-color: #343a40;
  border-color: #343a40;
}

.btn-dark:hover {
  color: #fff;
  background-color: #23272b;
  border-color: #1d2124;
}

.btn-outline-primary {
  color: #007bff;
  border-color: #007bff;
}

.btn-outline-primary:hover {
  color: #fff;
  background-color: #007bff;
  border-color: #007bff;
}

.btn-outline-secondary {
  color: #6c757d;
  border-color: #6c757d;
}

.btn-outline-secondary:hover {
  color: #fff;
  background-color: #6c757d;
  border-color: #6c757d;
}

.btn-outline-success {
  color: #28a745;
  border-color: #28a745;
}

.btn-outline-success:hover {
  color: #fff;
  background-color: #28a745;
  border-color: #28a745;
}

.btn-outline-danger {
  color: #dc3545;
  border-color: #dc3545;
}

.btn-outline-danger:hover {
  color: #fff;
  background-color: #dc3545;
  border-color: #dc3545;
}

.btn-outline-warning {
  color: #ffc107;
  border-color: #ffc107;
}

.btn-outline-warning:hover {
  color: #212529;
  background-color: #ffc107;
  border-color: #ffc107;
}

.btn-outline-info {
  color: #17a2b8;
  border-color: #17a2b8;
}

.btn-outline-info:hover {
  color: #fff;
  background-color: #17a2b8;
  border-color: #17a2b8;
}

.btn-outline-light {
  color: #f8f9fa;
  border-color: #f8f9fa;
}

.btn-outline-light:hover {
  color: #212529;
  background-color: #f8f9fa;
  border-color: #f8f9fa;
}

.btn-outline-dark {
  color: #343a40;
  border-color: #343a40;
}

.btn-outline-dark:hover {
  color: #fff;
  background-color: #343a40;
  border-color: #343a40;
}

.btn-sm {
  padding: 0.25rem 0.5rem;
  font-size: 0.875rem;
  line-height: 1.5;
  border-radius: 0.2rem;
}

.btn-lg {
  padding: 0.5rem 1rem;
  font-size: 1.25rem;
  line-height: 1.5;
  border-radius: 0.3rem;
}

.btn-block {
  display: block;
  width: 100%;
}

.card {
  position: relative;
  display: flex;
  flex-direction: column;
  min-width: 0;
  word-wrap: break-word;
  background-color: #fff;
  background-clip: border-box;
  border: 1px solid rgba(0, 0, 0, 0.125);
  border-radius: 0.25rem;
}

.card-body {
  flex: 1 1 auto;
  min-height: 1px;
  padding: 1.25rem;
}

.card-title {
  margin-bottom: 0.75rem;
}

.card-subtitle {
  margin-top: -0.375rem;
  margin-bottom: 0;
}

.card-text:last-child {
  margin-bottom: 0;
}

.card-header {
  padding: 0.75rem 1.25rem;
  margin-bottom: 0;
  background-color: rgba(0, 0, 0, 0.03);
  border-bottom: 1px solid rgba(0, 0, 0, 0.125);
}

.card-footer {
  padding: 0.75rem 1.25rem;
  background-color: rgba(0, 0, 0, 0.03);
  border-top: 1px solid rgba(0, 0, 0, 0.125);
}

.alert {
  position: relative;
  padding: 0.75rem 1.25rem;
  margin-bottom: 1rem;
  border: 1px solid transparent;
  border-radius: 0.25rem;
}

.alert-primary {
  color: #004085;
  background-color: #cce5ff;
  border-color: #b8daff;
}

.alert-secondary {
  color: #383d41;
  background-color: #e2e3e5;
  border-color: #d6d8db;
}

.alert-success {
  color: #155724;
  background-color: #d4edda;
  border-color: #c3e6cb;
}

.alert-danger {
  color: #721c24;
  background-color: #f8d7da;
  border-color: #f5c6cb;
}

.alert-warning {
  color: #856404;
  background-color: #fff3cd;
  border-color: #ffeeba;
}

.alert-info {
  color: #0c5460;
  background-color: #d1ecf1;
  border-color: #bee5eb;
}

.alert-light {
  color: #818182;
  background-color: #fefefe;
  border-color: #fdfdfe;
}

.alert-dark {
  color: #1b1e21;
  background-color: #d6d8d9;
  border-color: #c6c8ca;
}

.badge {
  display: inline-block;
  padding: 0.25em 0.4em;
  font-size: 75%;
  font-weight: 700;
  line-height: 1;
  text-align: center;
  white-space: nowrap;
  vertical-align: baseline;
  border-radius: 0.25rem;
  transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}

.badge-primary {
  color: #fff;
  background-color: #007bff;
}

.badge-secondary {
  color: #fff;
  background-color: #6c757d;
}

.badge-success {
  color: #fff;
  background-color: #28a745;
}

.badge-danger {
  color: #fff;
  background-color: #dc3545;
}

.badge-warning {
  color: #212529;
  background-color: #ffc107;
}

.badge-info {
  color: #fff;
  background-color: #17a2b8;
}

.badge-light {
  color: #212529;
  background-color: #f8f9fa;
}

.badge-dark {
  color: #fff;
  background-color: #343a40;
}

.form-control {
  display: block;
  width: 100%;
  height: calc(1.5em + 0.75rem + 2px);
  padding: 0.375rem 0.75rem;
  font-size: 1rem;
  font-weight: 400;
  line-height: 1.5;
  color: #495057;
  background-color: #fff;
  background-clip: padding-box;
  border: 1px solid #ced4da;
  border-radius: 0.25rem;
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}

.form-control:focus {
  color: #495057;
  background-color: #fff;
  border-color: #80bdff;
  outline: 0;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

.form-group {
  margin-bottom: 1rem;
}

.form-text {
  display: block;
  margin-top: 0.25rem;
}

.form-row {
  display: flex;
  flex-wrap: wrap;
  margin-right: -5px;
  margin-left: -5px;
}

.form-row > .col,
.form-row > [class*="col-"] {
  padding-right: 5px;
  padding-left: 5px;
}

.form-check {
  position: relative;
  display: block;
  padding-left: 1.25rem;
}

.form-check-input {
  position: absolute;
  margin-top: 0.3rem;
  margin-left: -1.25rem;
}

.form-check-label {
  margin-bottom: 0;
}

.form-check-inline {
  display: inline-flex;
  align-items: center;
  padding-left: 0;
  margin-right: 0.75rem;
}

.form-check-inline .form-check-input {
  position: static;
  margin-top: 0;
  margin-right: 0.3125rem;
  margin-left: 0;
}

.valid-feedback {
  display: none;
  width: 100%;
  margin-top: 0.25rem;
  font-size: 80%;
  color: #28a745;
}

.invalid-feedback {
  display: none;
  width: 100%;
  margin-top: 0.25rem;
  font-size: 80%;
  color: #dc3545;
}

.was-validated .form-control:valid, .form-control.is-valid {
  border-color: #28a745;
  padding-right: calc(1.5em + 0.75rem);
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' width='8' height='8' viewBox='0 0 8 8'%3e%3cpath fill='%2328a745' d='M2.3 6.73L.6 4.53c-.4-1.04.46-1.4 1.1-.8l1.1 1.4 3.4-3.8c.6-.63 1.6-.27 1.2.7l-4 4.6c-.43.5-.8.4-1.1.1z'/%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right calc(0.375em + 0.1875rem) center;
  background-size: calc(0.75em + 0.375rem) calc(0.75em + 0.375rem);
}

.was-validated .form-control:invalid, .form-control.is-invalid {
  border-color: #dc3545;
  padding-right: calc(1.5em + 0.75rem);
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' fill='none' stroke='%23dc3545' viewBox='0 0 12 12'%3e%3ccircle cx='6' cy='6' r='4.5'/%3e%3cpath stroke-linejoin='round' d='M5.8 3.6h.4L6 6.5z'/%3e%3ccircle cx='6' cy='8.2' r='.6' fill='%23dc3545' stroke='none'/%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right calc(0.375em + 0.1875rem) center;
  background-size: calc(0.75em + 0.375rem) calc(0.75em + 0.375rem);
}

.nav {
  display: flex;
  flex-wrap: wrap;
  padding-left: 0;
  margin-bottom: 0;
  list-style: none;
}

.nav-link {
  display: block;
  padding: 0.5rem 1rem;
}

.nav-link:hover, .nav-link:focus {
  text-decoration: none;
}

.nav-link.disabled {
  color: #6c757d;
  pointer-events: none;
  cursor: default;
}

.nav-tabs {
  border-bottom: 1px solid #dee2e6;
}

.nav-tabs .nav-item {
  margin-bottom: -1px;
}

.nav-tabs .nav-link {
  border: 1px solid transparent;
  border-top-left-radius: 0.25rem;
  border-top-right-radius: 0.25rem;
}

.nav-tabs .nav-link:hover, .nav-tabs .nav-link:focus {
  border-color: #e9ecef #e9ecef #dee2e6;
}

.nav-tabs .nav-link.disabled {
  color: #6c757d;
  background-color: transparent;
  border-color: transparent;
}

.nav-tabs .nav-link.active,
.nav-tabs .nav-item.show .nav-link {
  color: #495057;
  background-color: #fff;
  border-color: #dee2e6 #dee2e6 #fff;
}

.nav-pills .nav-link {
  border-radius: 0.25rem;
}

.nav-pills .nav-link.active,
.nav-pills .show > .nav-link {
  color: #fff;
  background-color: #007bff;
}

.navbar {
  position: relative;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1rem;
}

.navbar-brand {
  display: inline-block;
  padding-top: 0.3125rem;
  padding-bottom: 0.3125rem;
  margin-right: 1rem;
  font-size: 1.25rem;
  line-height: inherit;
  white-space: nowrap;
}

.navbar-brand:hover, .navbar-brand:focus {
  text-decoration: none;
}

.navbar-nav {
  display: flex;
  flex-direction: column;
  padding-left: 0;
  margin-bottom: 0;
  list-style: none;
}

.navbar-nav .nav-link {
  padding-right: 0;
  padding-left: 0;
}

.navbar-nav .dropdown-menu {
  position: static;
  float: none;
}

.navbar-text {
  display: inline-block;
  padding-top: 0.5rem;
  padding-bottom: 0.5rem;
}

.navbar-collapse {
  flex-basis: 100%;
  flex-grow: 1;
  align-items: center;
}

.navbar-toggler {
  padding: 0.25rem 0.75rem;
  font-size: 1.25rem;
  line-height: 1;
  background-color: transparent;
  border: 1px solid transparent;
  border-radius: 0.25rem;
}

.navbar-toggler:hover, .navbar-toggler:focus {
  text-decoration: none;
}

.navbar-toggler-icon {
  display: inline-block;
  width: 1.5em;
  height: 1.5em;
  vertical-align: middle;
  content: "";
  background: no-repeat center center;
  background-size: 100% 100%;
}

.navbar-light .navbar-brand {
  color: rgba(0, 0, 0, 0.9);
}

.navbar-light .navbar-brand:hover, .navbar-light .navbar-brand:focus {
  color: rgba(0, 0, 0, 0.9);
}

.navbar-light .navbar-nav .nav-link {
  color: rgba(0, 0, 0, 0.5);
}

.navbar-light .navbar-nav .nav-link:hover, .navbar-light .navbar-nav .nav-link:focus {
  color: rgba(0, 0, 0, 0.7);
}

.navbar-light .navbar-nav .nav-link.disabled {
  color: rgba(0, 0, 0, 0.3);
}

.navbar-light .navbar-nav .show > .nav-link,
.navbar-light .navbar-nav .active > .nav-link,
.navbar-light .navbar-nav .nav-link.show,
.navbar-light .navbar-nav .nav-link.active {
  color: rgba(0, 0, 0, 0.9);
}

.navbar-light .navbar-toggler {
  color: rgba(0, 0, 0, 0.5);
  border-color: rgba(0, 0, 0, 0.1);
}

.navbar-light .navbar-toggler-icon {
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' width='30' height='30' viewBox='0 0 30 30'%3e%3cpath stroke='rgba(0, 0, 0, 0.5)' stroke-linecap='round' stroke-miterlimit='10' stroke-width='2' d='M4 7h22M4 15h22M4 23h22'/%3e%3c/svg%3e");
}

.navbar-light .navbar-text {
  color: rgba(0, 0, 0, 0.5);
}

.navbar-light .navbar-text a {
  color: rgba(0, 0, 0, 0.9);
}

.navbar-light .navbar-text a:hover, .navbar-light .navbar-text a:focus {
  color: rgba(0, 0, 0, 0.9);
}

.navbar-dark .navbar-brand {
  color: #fff;
}

.navbar-dark .navbar-brand:hover, .navbar-dark .navbar-brand:focus {
  color: #fff;
}

.navbar-dark .navbar-nav .nav-link {
  color: rgba(255, 255, 255, 0.5);
}

.navbar-dark .navbar-nav .nav-link:hover, .navbar-dark .navbar-nav .nav-link:focus {
  color: rgba(255, 255, 255, 0.75);
}

.navbar-dark .navbar-nav .nav-link.disabled {
  color: rgba(255, 255, 255, 0.25);
}

.navbar-dark .navbar-nav .show > .nav-link,
.navbar-dark .navbar-nav .active > .nav-link,
.navbar-dark .navbar-nav .nav-link.show,
.navbar-dark .navbar-nav .nav-link.active {
  color: #fff;
}

.navbar-dark .navbar-toggler {
  color: rgba(255, 255, 255, 0.5);
  border-color: rgba(255, 255, 255, 0.1);
}

.navbar-dark .navbar-toggler-icon {
  background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' width='30' height='30' viewBox='0 0 30 30'%3e%3cpath stroke='rgba(255, 255, 255, 0.5)' stroke-linecap='round' stroke-miterlimit='10' stroke-width='2' d='M4 7h22M4 15h22M4 23h22'/%3e%3c/svg%3e");
}

.navbar-dark .navbar-text {
  color: rgba(255, 255, 255, 0.5);
}

.navbar-dark .navbar-text a {
  color: #fff;
}

.navbar-dark .navbar-text a:hover, .navbar-dark .navbar-text a:focus {
  color: #fff;
}

@media (min-width: 576px) {
  .navbar-expand-sm {
    flex-flow: row nowrap;
    justify-content: flex-start;
  }
  .navbar-expand-sm .navbar-nav {
    flex-direction: row;
  }
  .navbar-expand-sm .navbar-nav .dropdown-menu {
    position: absolute;
  }
  .navbar-expand-sm .navbar-nav .nav-link {
    padding-right: 0.5rem;
    padding-left: 0.5rem;
  }
  .navbar-expand-sm .navbar-collapse {
    display: flex !important;
    flex-basis: auto;
  }
  .navbar-expand-sm .navbar-toggler {
    display: none;
  }
}

@media (min-width: 768px) {
  .navbar-expand-md {
    flex-flow: row nowrap;
    justify-content: flex-start;
  }
  .navbar-expand-md .navbar-nav {
    flex-direction: row;
  }
  .navbar-expand-md .navbar-nav .dropdown-menu {
    position: absolute;
  }
  .navbar-expand-md .navbar-nav .nav-link {
    padding-right: 0.5rem;
    padding-left: 0.5rem;
  }
  .navbar-expand-md .navbar-collapse {
    display: flex !important;
    flex-basis: auto;
  }
  .navbar-expand-md .navbar-toggler {
    display: none;
  }
}

@media (min-width: 992px) {
  .navbar-expand-lg {
    flex-flow: row nowrap;
    justify-content: flex-start;
  }
  .navbar-expand-lg .navbar-nav {
    flex-direction: row;
  }
  .navbar-expand-lg .navbar-nav .dropdown-menu {
    position: absolute;
  }
  .navbar-expand-lg .navbar-nav .nav-link {
    padding-right: 0.5rem;
    padding-left: 0.5rem;
  }
  .navbar-expand-lg .navbar-collapse {
    display: flex !important;
    flex-basis: auto;
  }
  .navbar-expand-lg .navbar-toggler {
    display: none;
  }
}

@media (min-width: 1200px) {
  .navbar-expand-xl {
    flex-flow: row nowrap;
    justify-content: flex-start;
  }
  .navbar-expand-xl .navbar-nav {
    flex-direction: row;
  }
  .navbar-expand-xl .navbar-nav .dropdown-menu {
    position: absolute;
  }
  .navbar-expand-xl .navbar-nav .nav-link {
    padding-right: 0.5rem;
    padding-left: 0.5rem;
  }
  .navbar-expand-xl .navbar-collapse {
    display: flex !important;
    flex-basis: auto;
  }
  .navbar-expand-xl .navbar-toggler {
    display: none;
  }
}

.navbar-expand {
  flex-flow: row nowrap;
  justify-content: flex-start;
}

.navbar-expand .navbar-nav {
  flex-direction: row;
}

.navbar-expand .navbar-nav .dropdown-menu {
  position: absolute;
}

.navbar-expand .navbar-nav .nav-link {
  padding-right: 0.5rem;
  padding-left: 0.5rem;
}

.navbar-expand .navbar-collapse {
  display: flex !important;
  flex-basis: auto;
}

.navbar-expand .navbar-toggler {
  display: none;
}
`
}

// GenerateFullCSS generates the full CSS for web
func (a *WebAdapter) GenerateFullCSS() string {
	var css strings.Builder

	css.WriteString(a.GenerateResetCSS())
	css.WriteString(a.GenerateBaseCSS())
	css.WriteString(a.GenerateUtilitiesCSS())
	css.WriteString(a.GenerateComponentsCSS())

	return css.String()
}
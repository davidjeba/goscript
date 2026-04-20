# Vibe - Go Motion Library

`vibe` is the GoScript answer to Motion / Framer Motion style UI animation.

It is not a clone of the React package surface. The goal is to steal the capability envelope and express it in a Go-native way for GoScript and GoUIX.

## What Vibe should cover

- Declarative `initial`, `animate`, and `exit` targets
- Variants for coordinated parent-child motion
- Gesture-driven states like hover, tap, drag, focus, and in-view
- Motion values and derived values for scroll-linked interfaces
- Presence tracking for exit animations
- Transform-based layout animation and shared element transitions

## First-pass foundations now in code

- `pkg/vibe/transition.go`: animation timing and spring/tween/inertia transitions
- `pkg/vibe/motion_value.go`: subscribable motion values with velocity tracking
- `pkg/vibe/variant.go`: style targets, variants, gesture targets, and motion props
- `pkg/vibe/layout.go`: transform-style layout delta calculation
- `pkg/vibe/presence.go`: AnimatePresence-style lifecycle tracking
- `pkg/vibe/viewport.go`: scroll progress and in-view configuration

## Capability mapping

| Motion capability | Vibe direction |
| --- | --- |
| `motion.div` animation props | `MotionProps` + style targets |
| `variants` | `VariantSet` |
| `AnimatePresence` | `PresenceController` |
| `layout` / `layoutId` | `LayoutSnapshot` + `ComputeLayoutDelta` |
| `useScroll` / motion values | `MotionValue` + `ScrollState` |
| `useInView` / `whileInView` | `InViewOptions` + `WhileInView` |

## What is still missing

- Binding Vibe to real GoUIX and GoScript rendering
- Browser/runtime execution for transitions and keyframes
- Scroll observers and viewport tracking at runtime
- Gesture-to-variant execution in actual components
- Shared element transitions across route/view changes

## Why this matters

GoScript needs app-grade motion, not just CSS transitions. Products like docs, editors, ERPs, boards, dashboards, and social interfaces all need micro-interactions, layout transitions, drag feedback, and scroll-linked motion to feel native.

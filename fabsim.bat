@echo off
rem Convenience wrapper: run the fabsim CLI from the repo root via `go run`.
rem %~dp0 expands to the directory this .bat lives in, so invocation from a
rem subdirectory still resolves the cmd/fabsim package correctly.
rem
rem -pgo=auto applies cmd/fabsim/default.pgo (the captured anneal profile)
rem at build time. Go 1.21+ defaults to -pgo=auto already; naming it here
rem documents the intent.
go run -pgo=auto "%~dp0cmd\fabsim" %*

#!/usr/bin/env nu

const output_path = "bin"

def build-all [build_options: list<string>] {
    print "building all binaries..."

    let architectures = ["amd64", "arm64"]
    let operating_systems = ["windows", "linux", "darwin"]

    for arch in $architectures {
        for os in $operating_systems {
            $env.GOOS = $os
            $env.GOARCH = $arch

            mut binary_path = $"($output_path)/otf-($os)-($arch)"

            if $os == "windows" {
                $binary_path = $"($binary_path).exe"
            }

            let options = ["build", "-o", $binary_path] ++ $build_options
            run-external go ...$options
            print $"built ($binary_path)"
        }
    }
}

def build-local [build_options: list<string>] {
    print "building local binary..."
    let binary_path = $"($output_path)/otf"
    let options = ["build", "-o", $binary_path] ++ $build_options
    run-external go ...$options
    print $"built ($binary_path)"
}

def get-ldflags [version: string] {
    const go_import_path = "github.com/restechnica/opinionated-terraform"

    let ldflags = $"-X ($go_import_path)/internal/ldflags.Version=($version)"

    return $ldflags
}

def "main build" [--all, --version: string = "dev"] {
    let ldflags = get-ldflags $version

    let options: list<string> = ["-ldflags", $ldflags]

    if $all {
        build-all $options
    } else {
        build-local $options
    }
}

def main [] {}

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

OPTS = [
    "-DMCL_USE_VINT",
    "-DMCL_DONT_USE_OPENSSL",
    "-DMCL_LLVM_BMI2=0",
    "-DMCL_USE_LLVM=1",
    "-DMCL_VINT_FIXED_BUFFER",
    "-DMCL_SIZEOF_UNIT=8",
    "-DMCL_MAX_BIT_SIZE=384",
    "-DCYBOZU_DONT_USE_EXCEPTION",
    "-DCYBOZU_DONT_USE_STRING",
    "-DBLS_SWAP_G",
    "-DBLS_ETH",
]

genrule(
    name = "base64_ll",
    outs = ["src/base64.ll"],  # llvm assembly language file.
    tools = [
        "@herumi_mcl//:src_gen",
    ],
    cmd = "touch func.list && $(location @herumi_mcl//:src_gen) -u 64 -f func.list > $@",
)

genrule(
    name = "base64_o",
    srcs = [
        "src/base64.ll",
    ],
    outs = ["base64.o"],
    cmd = "external/llvm_toolchain/bin/clang++ -c -o $@ $(location src/base64.ll)",
    tools = ["@llvm_toolchain//:clang"],
)

cc_library(
    name = "lib",
    srcs = [
        "@herumi_mcl//:src/fp.cpp",
        "@herumi_bls//:src/bls_c384_256.cpp",
        "@herumi_bls//:src/bls_c_impl.hpp",
        ":base64_o",
    ],
    deps = ["@herumi_mcl//:bn"],
    includes = [
        "bls/include",
    ],
    hdrs = [
        "bls/include/bls/bls.h",
        "bls/include/bls/bls384_256.h",
        "bls/include/mcl/bn.h",
        "bls/include/mcl/bn_c384_256.h",
        "@herumi_mcl//:include/mcl/curve_type.h",
    ],
    copts = OPTS + [
        "-std=c++03",
    ],
    visibility = ["//visibility:public"],
)

cc_library(
    name = "precompiled",
    srcs = select({
        "@io_bazel_rules_go//go/platform:android_arm": [
            "bls/lib/android/armeabi-v7a/libbls384_256.a",
        ],
        "@io_bazel_rules_go//go/platform:linux_arm64": [
            "bls/lib/android/arm64-v8a/libbls384_256.a",
        ],
        "@io_bazel_rules_go//go/platform:android_arm64": [
            "bls/lib/android/arm64-v8a/libbls384_256.a",
        ],
        "@io_bazel_rules_go//go/platform:darwin_amd64": [
            "bls/lib/darwin/amd64/libbls384_256.a",
        ],
        "@io_bazel_rules_go//go/platform:linux_amd64": [
            "bls/lib/linux/amd64/libbls384_256.a",
        ],
        "@io_bazel_rules_go//go/platform:windows_amd64": [
            "bls/lib/windows/amd64/libbls384_256.a",
        ],
        "//conditions:default": [],
    }),
    hdrs = [
        "bls/include/bls/bls.h",
        "bls/include/bls/bls384_256.h",
        "bls/include/mcl/bn.h",
        "bls/include/mcl/bn_c384_256.h",
        "bls/include/mcl/curve_type.h",
    ],
    includes = [
        "bls/include",
    ],
    deprecation = "Using precompiled BLS archives. To build BLS from source with llvm, use --config=llvm.",
)

config_setting(
    name = "llvm_compiler_enabled",
    define_values = {
        "compiler": "llvm",
    },
)

go_library(
    name = "go_default_library",
    importpath = "github.com/herumi/bls-eth-go-binary/bls",
    srcs = [
        "bls/bls.go",
        "bls/callback.go",
        "bls/cast.go",
        "bls/mcl.go",
    ],
    cdeps = select({
        ":llvm_compiler_enabled": [":lib"],
        "//conditions:default": [":precompiled"],
    }),
    cgo = True,
    visibility = ["//visibility:public"],
    clinkopts = [
        "-Wl,--unresolved-symbols=ignore-all",  # Ignore missing asan symbols
    ],
)

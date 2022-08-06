workspace "RDIPs-Broker"
    configurations {"Debug", "Release"}
    location "build"
    architecture "x64"

outputDir = "%{cfg.buildcfg}-%{cfg.system}-%{cfg.architecture}"

Includedirs = {}
Includedirs["amqcpp"] =  "./dependencies/AMQP-CPP/include"


include "./dependencies"

project "RDIPs-Broker"
    kind "ConsoleApp"
    language "C++"
    targetdir("./bin/%{prj.name}")
    objdir("./bin-int/%{prj.name}")
    files {"./src/**.h", "./src/**.cpp"}
    
    includedirs {
        "%{Includedirs.amqcpp}",
        "./include",
        "./src"
    }

    links {
        "amqcpp",
        "pthread",
        "ev",
        "dl"
    }

    filter "configurations:Debug"
        defines {"DEBUG"}
        symbols "On"

    filter "configurations:Release"
        defines {"NDEBUG"}
        optimize "On"
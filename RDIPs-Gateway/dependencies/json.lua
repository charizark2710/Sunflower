project "json"
    kind "StaticLib"
    language "C++"
    pic "on"
    targetdir("./bin/%{prj.name}")
    cppdialect "c++17"
    objdir("./bin-int/%{prj.name}")
    files {"./json/include/**.hpp", "./json/single_include/**.hpp"}
    
    includedirs {
        "./json/include",
        "./json/single_include"
    }
    
    filter "configurations:Debug"
        defines {"DEBUG"}
        symbols "On"
    
    filter "configurations:Release"
        defines {"NDEBUG"}
        optimize "On"
    
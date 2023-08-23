#!/bin/bash

set -e

ProductName=gtbox
OSTYPE="Unknown"
GetOSType(){
    uNames=`uname -s`
    osName=${uNames: 0: 4}
    if [ "$osName" == "Darw" ] # Darwin
    then
        OSTYPE="Darwin"
    elif [ "$osName" == "Linu" ] # Linux
    then
        OSTYPE="Linux"
    elif [ "$osName" == "MING" ] # MINGW, windows, git-bash
    then
        OSTYPE="Windows"
    else
        OSTYPE="Unknown"
    fi
}
GetOSType

removeCache() {
    rm -rf ./${ProductName}_config.go
    rm -rf ./install_${ProductName}.sh
}

install() {
    echo ${OSTYPE}

    complate_gopath_dir=${GOPATH}
    if [ ${OSTYPE} == "Windows" ]
    then
        ago_path_dir=`echo "${GOPATH/':\\'/'/'}" | sed 's/\"//g'`
        complate_gopath_dir='/'`echo "${ago_path_dir}" | tr A-Z a-z`
    fi

    find ${complate_gopath_dir}/pkg/mod/github.com/george012 -depth -name "${ProductName}@*" -exec rm -rf {} \;

    go get -u github.com/george012/${ProductName}@latest

    wget --no-check-certificate https://raw.githubusercontent.com/george012/${ProductName}/master/config/config.go -O ${ProductName}_config.go \
    && {

        aVersionNo=$(grep ProjectVersion ${ProductName}_config.go | awk -F '"' '{print $2}' | sed 's/\"//g') \
        && CustomLibs=$(ls -l ${complate_gopath_dir}/pkg/mod/github.com/george012/gtbox@v$aVersionNo/libs |awk '/^d/ {print $NF}') \
        && for alibName in ${CustomLibs}
        do
            if [ ${OSTYPE} == "Darwin" ]; then # Darwin
                srcPWD=`pwd`
        #        cd ${GOPATH}/pkg/mod/github.com/george012/gtbox@v${aVersionNo} && /Applications/Xcode.app/Contents/Developer/Toolchains/XcodeDefault.xctoolchain/usr/bin/install_name_tool -add_rpath ../gtbox@v${aVersionNo} ${produckName} && cd ${srcPWD}
                sudo ln -s ${complate_gopath_dir}/pkg/mod/github.com/george012/${ProductName}@v${aVersionNo}/libs/${alibName}/lib${alibName}.dylib /usr/local/lib/lib${alibName}.dylib
                sudo ln -s /usr/local/lib/lib${alibName}.dylib /usr/local/lib/lib${alibName}_arm64.dylib
            elif [ ${OSTYPE} == "Linux" ] # Linux
            then
                ln -s ${complate_gopath_dir}/pkg/mod/github.com/george012/${ProductName}@v${aVersionNo}/libs/${alibName}/lib${alibName}.so /lib64/lib${alibName}.so && ldconfig
            elif [ ${OSTYPE} == "Windows" ] # MINGW, windows, git-bash
            then
                ln -s ${complate_gopath_dir}/pkg/mod/github.com/george012/${ProductName}@v${aVersionNo}/libs/${alibName}/${alibName}.dll /c/Windows/System32/${alibName}.dll
            else
                echo ${OSTYPE}
            fi
        done
    }

    removeCache
}

uninstall() {
    complate_gopath_dir=${GOPATH}

    # 找到所有版本的库并删除
    find ${complate_gopath_dir}/pkg/mod/github.com/george012/${ProductName}@* -type d -exec rm -rf {} \;

    # 删除所有自定义库
    CustomLibs=$(ls -l ${complate_gopath_dir}/pkg/mod/github.com/george012/${ProductName}/libs |awk '/^d/ {print $NF}')

    for libName in ${CustomLibs}
    do
        if [ ${OSTYPE} == "Darwin" ] # Darwin
        then
            rm -rf /usr/local/lib/lib${libName}_arm64.dylib
            rm -rf /usr/local/lib/lib${libName}.dylib
        elif [ ${OSTYPE} == "Linux" ] # Linux
        then
            rm -rf /lib64/lib${libName}.so
        elif [ ${OSTYPE} == "Windows" ] # MINGW, windows, git-bash
        then
            ago_path_dir=`echo "${GOPATH/':\\'/'/'}" | sed 's/\"//g'`
            complate_gopath_dir='/'`echo "${ago_path_dir}" | tr A-Z a-z`
            rm -rf /c/Windows/System32/${libName}.dll
        else
            echo ${OSTYPE}
        fi
    done

    removeCache
}

echo "============================ ${ProductName} ============================"
echo "  1、安装 ${ProductName}"
echo "  2、卸载 ${ProductName}"
echo "======================================================================"
read -p "$(echo -e "请选择[1-2]：")" choose
case $choose in
1)
    install
    ;;
2)
    uninstall
    ;;
*)
    echo "输入错误，请重新输入！"
    ;;
esac

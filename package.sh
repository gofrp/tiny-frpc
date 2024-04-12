#!/bin/sh
set -e

# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

tiny_frpc_version=`./bin/frpc-gssh --v`
echo "build version: $tiny_frpc_version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd android'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle riscv64'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        tiny_frpc_dir_name="tiny-frpc_${tiny_frpc_version}_${os}_${arch}"
        tiny_frpc_path="./packages/tiny-frpc_${tiny_frpc_version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./frpc-gssh_${os}_${arch}.exe" ]; then
                continue
            fi
            if [ ! -f "./frpc-nssh_${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${tiny_frpc_path}
            mv ./frpc-gssh_${os}_${arch}.exe ${tiny_frpc_path}/frpc-gssh.exe
            mv ./frpc-nssh_${os}_${arch}.exe ${tiny_frpc_path}/frpc-nssh.exe
        else
            if [ ! -f "./frpc-gssh_${os}_${arch}" ]; then
                continue
            fi
            if [ ! -f "./frpc-nssh_${os}_${arch}" ]; then
                continue
            fi
            mkdir ${tiny_frpc_path}
            mv ./frpc-gssh_${os}_${arch} ${tiny_frpc_path}/frpc-gssh
            mv ./frpc-nssh_${os}_${arch} ${tiny_frpc_path}/frpc-nssh
        fi  
        cp ../LICENSE ${tiny_frpc_path}
        cp -f ../conf/frpc.toml ${tiny_frpc_path}

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${tiny_frpc_dir_name}.zip ${tiny_frpc_dir_name}
        else
            tar -zcf ${tiny_frpc_dir_name}.tar.gz ${tiny_frpc_dir_name}
        fi  
        cd ..
        rm -rf ${tiny_frpc_path}
    done
done

cd -

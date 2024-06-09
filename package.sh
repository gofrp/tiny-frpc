#!/bin/sh
set -e

# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

tiny_frpc_version=`./bin/tiny-frpc --v`
echo "build version: $tiny_frpc_version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd android'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle riscv64'
extra_all='_ 6'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        for extra in $extra_all; do
            suffix="${os}_${arch}"
            if [ "x${extra}" != x"_" ]; then
                suffix="${os}_${arch}_${extra}"
            fi
            tiny_frpc_dir_name="tiny-frpc_${tiny_frpc_version}_${suffix}"
            tiny_frpc_path="./packages/tiny-frpc_${tiny_frpc_version}_${suffix}"

            if [ "x${os}" = x"windows" ]; then
                if [ ! -f "./tiny-frpc_${os}_${arch}.exe" ]; then
                    continue
                fi
                if [ ! -f "./tiny-frpc-ssh_${os}_${arch}.exe" ]; then
                    continue
                fi
                mkdir ${tiny_frpc_path}
                mv ./tiny-frpc_${os}_${arch}.exe ${tiny_frpc_path}/frpc.exe
                mv ./tiny-frpc-ssh_${os}_${arch}.exe ${tiny_frpc_path}/tiny-frpc-ssh.exe
            else
                if [ ! -f "./tiny-frpc_${suffix}" ]; then
                    continue
                fi
                if [ ! -f "./tiny-frpc-ssh_${suffix}" ]; then
                    continue
                fi
                mkdir ${tiny_frpc_path}
                mv ./tiny-frpc_${suffix} ${tiny_frpc_path}/tiny-frpc
                mv ./tiny-frpc-ssh_${suffix} ${tiny_frpc_path}/tiny-frpc-ssh
            fi
            cp ../LICENSE ${tiny_frpc_path}
            cp -f ../conf/frpc.toml ${tiny_frpc_path}

            # packages
            cd ./packages
            if [ "x${os}" = x"windows" ]; then
                zip -rq ${tiny_frpc_dir_name}.zip ${tiny_frpc_dir_name}
                echo "windows"
            else
                tar -zcf ${tiny_frpc_dir_name}.tar.gz ${tiny_frpc_dir_name}
                echo "linux+mac"
            fi
            cd ..
            rm -rf ${tiny_frpc_path}
        done
    done
done

cd -
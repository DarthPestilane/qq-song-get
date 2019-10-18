## 安装

- 通过 go:

```sh
$ go get -u -v github.com/DarthPestilane/qq-song-get
```

- 通过 docker:

```sh
$ docker pull darthminion/qq-song-get
```

## 使用

- 通过本地二进制文件：

```sh
$ ./qq-song-get --color=off https://y.qq.com/n/yqq/album/000dilOO3JYIr4.html
```

- 通过 docker:

```sh
$ docker run --rm -it -v `pwd`/downloads:/downloads darthminion/qq-song-get https://y.qq.com/n/yqq/album/000dilOO3JYIr4.html
```

## 致谢

- [winterssy/music-get](https://github.com/winterssy/music-get)

## 免责声明

- 本项目仅供学习研究使用，禁止商业用途。
- 本项目使用的接口如无特别说明均为官方接口，音乐版权归源音乐平台所有，侵删。

## License

GPLv3

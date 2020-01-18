# connected [![GitHub release](https://img.shields.io/github/release/k1LoW/connected.svg)](https://github.com/k1LoW/connected/releases)

:electric_plug: Watch your MacBook connection :zap:

## Usage

**Power cable:**

``` console
$ connected watch -- say "Power cable disconnected."
```

or use `--command (-c)`

``` console
$ connected watch -c "sh ./slack_notify.sh"
```

**Wi-Fi:**

``` console
$ connected watch --wifi -- say "Wi-Fi disconnected."
```

## Install

**homebrew tap:**

```console
$ brew install k1LoW/tap/connected
```

**manually:**

Download binany from [releases page](https://github.com/k1LoW/connected/releases)

**go get:**

```console
$ go get github.com/k1LoW/connected
```

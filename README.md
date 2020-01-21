# connected [![GitHub release](https://img.shields.io/github/release/k1LoW/connected.svg)](https://github.com/k1LoW/connected/releases)

:electric_plug: Watch your MacBook connection :zap:

## Usage

**Watch power cable connection:**

``` console
$ connected watch -- say "Power cable disconnected."
```

or use `--command (-c)`

``` console
$ connected watch -c "osascript -e "set Volume 10"; say -v Alex "Power cable disconnected."
```

**Watch Wi-Fi connection:**

``` console
$ connected watch --wifi -- say "Wi-Fi disconnected."
```

**Watch Bluetooth devices connection:**

``` console
$ connected watch --bluetooth -- sh ./twilio_call.sh
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

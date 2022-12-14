
export default function (timespan: number) {
  if (timespan < 0) {
    return (<span>-</span>)
  }
  let time = timespan / 1e6;
  if (time < 1000) {
    return (<span>{time.toFixed(1)}ms</span>);
  }
  time = time / 1e3;
  if (time < 120) {
    return (<span>{time.toFixed(1)}s</span>);
  }
  let h = time / 3600;
  let m = (time % 3600) / 60;
  let s = time % 60;
  return (<span>{h}:{m}:{s}</span>)
}
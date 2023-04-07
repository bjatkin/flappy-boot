import { useRef, useState } from 'react'
import './App.css'
import { setupGo } from './components/wasm_exec'
import { setupTinyGo } from './components/tiny_wasm_exec'
import { init, parse_rom_header, Emulator } from '@uyouii/rustboyadvance-wasm'

function App() {
  let [loading, setLoading] = useState(true)
  let [running, setRunning] = useState(false)

  const middleTopRef = useRef<HTMLDivElement>(null)
  const middleBottomRef = useRef<HTMLDivElement>(null)
  const appRef = useRef<HTMLDivElement>(null)
  const runRef = useRef<HTMLButtonElement>(null)
  const overlayRef = useRef<HTMLDivElement>(null)

  init()

  let romData = new Uint8Array;
  let biosData = new Uint8Array;
  fetch("flappy_boot.gba").then((result) => {
    result.arrayBuffer().then((data) => {
      romData = new Uint8Array(data)
      let rom_info = parse_rom_header(romData);

      console.log("Game Code" + rom_info.get_game_code());
      console.log("Game Title" + rom_info.get_game_title());
      console.log("HERE 1")

      let emulator = new Emulator(biosData, romData);
      emulator.skip_bios()
      console.log("HERE 2")

      let canvas = document.getElementById("screen") as HTMLCanvasElement;
      if (canvas == null) {
        console.log("canvas was null")
        return
      }
      let ctx = canvas.getContext('2d') as CanvasRenderingContext2D;

      setInterval(() => {
        emulator.run_frame(ctx);
        console.log("HERE 3")
      }, 16)

    })
  })

  // use this line for wasm files complied with the go complier
  // setupGo();

  // use this line for wasm files complie with the tinygo complier
  // setupTinyGo();

  // @ts-ignore this is added to the global object by the setupGo/ setupTinyGo func
  // const go = new Go();

  // let mod: WebAssembly.Module
  // let inst: WebAssembly.Instance
  // WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
  //   mod = result.module;
  //   inst = result.instance;
    
  //   // loading is finished
  //   setLoading(false)
  //   if (appRef.current != null) {
  //     appRef.current.style.setProperty("--bg-1", "#833AB4")
  //     appRef.current.style.setProperty("--bg-2", "#FD1D1D")
  //     appRef.current.style.setProperty("--bg-3", "#FCB045")
  //     appRef.current.style.setProperty("--speed", "10s") 
  //   }

  //   if (runRef.current != null) {
  //     runRef.current.style.display = "unset"
  //   }
  // }).catch((err) => {
  //   console.error(err);
  // });

  async function run() {
    console.clear();
    // TODO: start the emulator here
    //
    // await go.run(inst);
    // inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
  }

  const handleClick = ():void => {
    if (overlayRef.current != null) {
      overlayRef.current.classList.add("fade-out")
    }

    window.setTimeout(() => {
      setRunning(true)

      // actually run the wasm code
      run();
    }, 500)
  }

  return (
  <div>{ running &&

    <div className='App' ref={appRef}>
      <div id="overlay" ref={overlayRef}></div>
      <div className='pannel'></div>
      <div id={loading ? "" : "middle-top"} className='pannel' ref={middleTopRef}></div>
      <div className='pannel'></div>

      <div className='pannel'></div>
      <div id={loading ? "" : "middle-bottom"} className='pannel' ref={middleBottomRef}></div>
      <div className='pannel'></div>

      <button id='runButton' ref={runRef} onClick={handleClick}>
        <svg id="runIcon" height="80" width="80" viewBox='0 0 100 100'>
          <polygon fill="#151515" points="35,25 75,50 35,75"></polygon>
          <circle id="runCircle" cx="50" cy="50" r="58" fill="#151515a0" stroke="#151515" strokeWidth="30"></circle>
        </svg>
      </button>
    </div>

  }</div>
  )
}

export default App
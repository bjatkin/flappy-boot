import { useRef, useState } from 'react'
import './App.css'
import { setupGo } from './components/wasm_exec'
import { setupTinyGo } from './components/tiny_wasm_exec'

function App() {
  let [loading, setLoading] = useState(true)
  let [running, setRunning] = useState(false)

  const middleTopRef = useRef<HTMLDivElement>(null)
  const middleBottomRef = useRef<HTMLDivElement>(null)
  const appRef = useRef<HTMLDivElement>(null)
  const runRef = useRef<HTMLButtonElement>(null)
  const overlayRef = useRef<HTMLDivElement>(null)

  setupGo();

  // @ts-ignore this is added to the global object by the setupGo/ setupTinyGo func
  const go = new Go();

  let mod: WebAssembly.Module
  let inst: WebAssembly.Instance
  WebAssembly.instantiateStreaming(fetch("flappy_boot.wasm"), go.importObject).then((result) => {
    mod = result.module;
    inst = result.instance;
    
    // loading is finished
    setLoading(false)
    if (appRef.current != null) {
      appRef.current.style.setProperty("--bg-1", "#833AB4")
      appRef.current.style.setProperty("--bg-2", "#FD1D1D")
      appRef.current.style.setProperty("--bg-3", "#FCB045")
      appRef.current.style.setProperty("--speed", "10s") 
    }

    if (runRef.current != null) {
      runRef.current.style.display = "unset"
    }
  }).catch((err) => {
    console.error(err);
  });

  async function run() {
    console.clear();
    await go.run(inst);
    inst = await WebAssembly.instantiate(mod, go.importObject); // reset instance
  }

  const handleClick = ():void => {
    if (overlayRef.current != null) {
      overlayRef.current.classList.add("fade-out")
    }

    window.setTimeout(() => {
      setRunning(true)

      // This prevents the bug where firefox builds are not full screen
      let root = document.getElementById("root")
      if (root != null) {
        root.classList.add('small');
      }

      run();
    }, 500)
  }

  return (
  <div className='game'>{ !running &&

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

import { useState } from 'react'
import './App.css'

function App() {


  return (
    <main className='flex h-dvh justify-center items-center bg-gray-100'>
      <div className='grid grid-cols-3 grid-rows-4 m-px rounded-xl outline-7 outline-[#e2dbc3] max-w-9/10 h-19/20 overflow-clip'>
        <img src="../public/GUIf2.png" className='flex grow col-span-full row-span-full h-dvh' />
        <div className='bg-[#848f67] row-span-full row-start-4 row-end-4 col-start-1 col-end-4 h-7/10 self-end'></div>
      </div>
    </main>
  )
}

export default App

import { useEffect, useState } from 'react'

import './App.css'
import { client } from './connect'
import UploadBookForm from './pages/uploads'
import BooksList from './pages/books'

function App() {

  useEffect( ()=>{
    client.getBooks({}).then((res)=>{
      console.log(res.books)
    })

  },[])
  return (
    <>
  <div>
<BooksList/>

<hr />
<hr />
<hr />
<hr />
<UploadBookForm/>
  </div>
    </>
  )
}

export default App

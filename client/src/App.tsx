import Navigation from './components/Navigation'
import './App.css'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Login from './pages/Login'
import { isLoggedIn } from './utils'
import NotFound from './pages/NotFound'
import UserProvider from './context/userContext'
import Groups from './pages/Groups'
import Movies from './pages/Movies'
import CreateGroup from './pages/CreateGroup'
import ViewGroup from './pages/ViewGroup'

function App() {
  return (
    <>
      <UserProvider>
        <BrowserRouter>
          <Navigation />
          <Routes>
            {isLoggedIn() && <Route path="/groups" element={<Groups />} /> }
            {isLoggedIn() && <Route path="/groups/create" element={<CreateGroup />} /> }
            {isLoggedIn() && <Route path="/groups/:groupId" element={<ViewGroup />} /> }
            {isLoggedIn() && <Route path="/movies" element={<Movies />} /> }
            <Route path="/*" element={<NotFound />} />
            <Route path="/" element={isLoggedIn() ? <Navigate to="/groups"/> : <Login />}/>
          </Routes>
        </BrowserRouter>
      </UserProvider>
    </>
  )
}

export default App

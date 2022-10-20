import "./Sidebar.scss";
import Sidebar from "./Components/organisms/common/Sidebar.organism";
import React from "react";
function Sidebar() {
    <>
    <Sidebar />
      <BrowserRouter>
      <Routes>
        <Route path='/' element={<TempPage />} />
      </Routes>
      </BrowserRouter>
    
    </>
  };
  

export default Sidebar;
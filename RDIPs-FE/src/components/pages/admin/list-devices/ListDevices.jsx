import TableAtom from '../../../atoms/table/Table.atom';
import ModalAtomProps from '../../../atoms/modal/Modal.atom'
import './ListDevices.scss'
import React from 'react';

const ListDevices = () => {
  const [openDetailModal, setOpenDetailModal] = React.useState(
    {status: true, 
      data : 
      {
        idDevice: "Frozen yoghurt",
        firmWareVer: 159,
        appVer: 6,
        common: 24,
        action: 4
      }});

  const handleOnCloseModal = () => setOpenDetailModal({...openDetailModal, status: false});

  const openDetailPopup = (detail) => {
    setOpenDetailModal({...openDetailModal, status: true, data: detail})
    console.log(detail);
  }

  return (
  <div className='list-devices-container'>
    <TableAtom onRowClick = {openDetailPopup} />
    <ModalAtomProps open={openDetailModal.status} handleClose = {handleOnCloseModal} data={openDetailModal.data}/>
  </div>
  )
}

export default ListDevices;
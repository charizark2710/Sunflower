import { HeadCell, ReceiptData } from '../../../../../../utils/interface';
import TableAtom from '../../../../../atoms/table/Table.atom';

const Receipt = () => {
  const receiptListData = [
    createData('R1', '2009-09-01T00:00:00.000Z', 'Component 1', 'Service 1', 10000),
    createData('R2', '2009-09-01T00:00:00.000Z', 'Component 3', 'Service 2', 10000),
    createData('R3', '2009-09-01T00:00:00.000Z', 'Component 4', 'Service 3', 10000),
    createData('R4', '2009-09-01T00:00:00.000Z', 'Component 5', 'Service 4', 10000),
    createData('R5', '2009-09-01T00:00:00.000Z', 'Component 6', 'Service 5', 10000),
  ];

  function createData(receive_id: string, datetime: string, component: string, service: string, total: number): ReceiptData {
    return {
      receive_id,
      datetime,
      component,
      service,
      total,
    };
  }

  const headCells: HeadCell[] = [
    {
      numeric: false,
      disablePadding: false,
      label: 'STT',
    },
    {
      id: 'receipt_id',
      numeric: false,
      disablePadding: false,
      label: 'Receipt ID',
    },
    {
      id: 'datetime',
      numeric: true,
      disablePadding: false,
      label: 'Datetime',
    },
    {
      id: 'component',
      numeric: true,
      disablePadding: false,
      label: 'Component',
    },
    {
      id: 'service',
      numeric: true,
      disablePadding: false,
      label: 'Service',
    },
    {
      id: 'total',
      numeric: true,
      disablePadding: false,
      label: 'Total',
    },
  ];

  function onRowClick() {
    return;
  }

  const receiptColumns = ['datetime', 'datetime', 'component', 'service', 'total'];
  return (
    <div>
      <TableAtom
        onRowClick = {onRowClick}
        rows={receiptListData}
        deviceColumns={receiptColumns}
        title='Receipt Table'
        headCells={headCells}
      />
    </div>
  );
};

export default Receipt;

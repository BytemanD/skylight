import {
    DomainTable, ProjectTable, RoleTable, UserTable, HypervisortTable,
    VolumeTypeTable, 
    FlavorDataTable,VolumeServiceTable ,KeypairDataTable, UsageTable, ComputeServiceTable,
    RegionTable, AZDataTable,
    SecurityGroupDataTable, QosPolicyDataTable, 
    ClusterTable,
    AggDataTable, MigrationDataTable, EndpointTable,
} from './tables.jsx';
import {
    VolumeDataTable, BackupDataTable, SnapshotDataTable,
    RouterDataTable, NetDataTable, PortDataTable,
    ServerDataTable,
} from './data_tables.js';

export const userTable = new UserTable();
export const projectTable = new ProjectTable();
export const roleTable = new RoleTable();
export const domainTable = new DomainTable();
export const endpointTable = new EndpointTable();

export const hypervisorTable = new HypervisortTable();
export const volumeTable = new VolumeDataTable();
export const volumeTypeTable = new VolumeTypeTable();
export const snapshotTable = new SnapshotDataTable();
export const flavorTable = new FlavorDataTable();
export const backupTable = new BackupDataTable();
export const volumeComputeServiceTable = new VolumeServiceTable();

export const keypairTable = new KeypairDataTable();
export const serverTable = new ServerDataTable();
export const usageTable = new UsageTable();
export const aggTable = new AggDataTable();
export const migrationTable = new MigrationDataTable();

export const serviceTable = new ComputeServiceTable();
export const routerTable = new RouterDataTable();
export const netTable = new NetDataTable();
export const portTable = new PortDataTable();
export const sgTable = new SecurityGroupDataTable();

export const qosPolicyTable = new QosPolicyDataTable();
export const clusterTable = new ClusterTable();
export const regionTable = new RegionTable();
export const azTable = new AZDataTable();

export default null;

// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {contract} from '../models';
import {context} from '../models';

export function AddRepo(arg1:contract.Repository):Promise<void>;

export function CheckPRs():Promise<Array<contract.Pullrequest>>;

export function OnShutdown(arg1:context.Context):Promise<void>;

export function Repos():Promise<Array<contract.Repository>>;

export function UpdateSchedule(arg1:string):Promise<void>;

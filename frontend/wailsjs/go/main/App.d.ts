// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function DeleteProcess(arg1:number):Promise<Error>;

export function GetLogs(arg1:number):Promise<Array<main.Log>>;

export function GetProcesses():Promise<Array<main.Process>>;

export function InsertProcess(arg1:main.Process):Promise<Error>;

export function OpenGithub():Promise<void>;

export function RunProcess(arg1:number):Promise<Error>;

export function StopProcess(arg1:number):Promise<Error>;

export function UpdateProcess(arg1:main.Process):Promise<Error>;

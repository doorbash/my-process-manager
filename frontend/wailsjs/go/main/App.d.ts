// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {main} from '../models';

export function DeleteLogs(arg1:number):Promise<void>;

export function DeleteProcess(arg1:number):Promise<void>;

export function GetLogs(arg1:number):Promise<Array<main.Log>>;

export function GetProcesses():Promise<Array<main.Process>>;

export function InsertProcess(arg1:main.Process):Promise<void>;

export function OpenGithub():Promise<void>;

export function ProcessesReorder(arg1:Array<number>):Promise<boolean>;

export function RunProcess(arg1:number):Promise<void>;

export function StopProcess(arg1:number):Promise<void>;

export function UpdateProcess(arg1:main.Process):Promise<void>;

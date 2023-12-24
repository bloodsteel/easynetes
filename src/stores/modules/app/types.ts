export interface AppState {
  theme: string;
  navbar: boolean;
  menu: boolean;
  topMenu: boolean;
  hideMenu: boolean;
  menuCollapse: boolean;
  footer: boolean;
  menuWidth: number;
  device: string;
  [key: string]: unknown;
}

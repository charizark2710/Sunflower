import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, ReplaySubject } from 'rxjs';
import { MenuItem } from '../interface/menu.model';
import { _ } from '@ngx-translate/core';

export const MENU_TITLE = {
  home: 'labels.menu.home',
  devices: 'labels.menu.listDevices',
  users: 'labels.menu.listUsers',
  admins: 'labels.menu.listAdmins',
  campaign: 'labels.menu.campaign',
  support: 'labels.menu.support',
  settings: 'labels.menu.settings',
};

export const ICON_URL = {
  home: '/assets/icons/home.svg',
  devices: '/assets/icons/list-devices.svg',
  users: '/assets/icons/list-users.svg',
  admins: '/assets/icons/list-admins.svg',
  campaign: '/assets/icons/campaign.svg',
  support: '/assets/icons/support.svg',
  settings: '/assets/icons/settings.svg',
};

export const URL_PATH = {
  home: '/',
  devices: '/devices',
  users: '/users',
  admins: '/admins',
  campaign: '/campaign',
  support: '/support',
  settings: '/settings',
};

export const MENU_ITEMS: MenuItem[] = [
  {
    name: MENU_TITLE.home,
    icon: ICON_URL.home,
    url: URL_PATH.home
  },
  {
    name: MENU_TITLE.devices,
    icon: ICON_URL.devices,
    url: URL_PATH.devices
  },
  {
    name: MENU_TITLE.users,
    icon: ICON_URL.users,
    url: URL_PATH.users
  },
  {
    name: MENU_TITLE.admins,
    icon: ICON_URL.admins,
    url: URL_PATH.admins
  },
  {
    name: MENU_TITLE.campaign,
    icon: ICON_URL.campaign,
    url: URL_PATH.campaign
  },
];

export const SUB_MENU_ITEMS: MenuItem[] = [
  {
    name: MENU_TITLE.support,
    icon: ICON_URL.support,
    url: URL_PATH.support
  },
  {
    name: MENU_TITLE.settings,
    icon: ICON_URL.settings,
    url: URL_PATH.settings
  },
];

@Injectable({
  providedIn: 'root',
})
export class MenuItemService {
  private _menuItems$ = new ReplaySubject<MenuItem[]>(1);
  private _subMenuItems$ = new ReplaySubject<MenuItem[]>(1);
  private _activeMenuItem?: MenuItem;

  constructor() {
    this._menuItems$.next(MENU_ITEMS);
    this._subMenuItems$.next(SUB_MENU_ITEMS);
  }

  get activeMenuItem() : MenuItem | undefined {
    return this._activeMenuItem;
  }

  get menuItems$(): Observable<MenuItem[]> {
    return this._menuItems$;
  }

  get subMenuItems$(): Observable<MenuItem[]> {
    return this._subMenuItems$;
  }

  setActiveMenuItem(menuItem?: MenuItem): void {
    this._activeMenuItem = menuItem;
  }
}

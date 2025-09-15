import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { combineLatest, Subject, take, takeUntil, tap } from 'rxjs';
import { MenuItem } from '../../interface/menu.model';
import { MenuItemService } from '../../services/menu-item.service';

@Component({
  selector: 'app-navigation',
  standalone: true,
  imports: [CommonModule, TranslateModule, RouterModule],
  templateUrl: './navigation.component.html',
  styleUrl: './navigation.component.scss',
})
export class NavigationComponent {
  destroy$ = new Subject<void>();
  menuItems!: MenuItem[];
  subMenuItems!: MenuItem[];
  allMenuItems: MenuItem[] = [];

  get activeMenuItem(): MenuItem | undefined {
    return this.menuItemService.activeMenuItem;
  }

  constructor(
    private menuItemService: MenuItemService,
    private router: Router
  ) {}

  ngOnInit() {
    this.loadMenuItems();
    this.subscribeToRouterEvents();
    this.updateActiveMenuItem();
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }

  setActiveMenuItem(menuItem: MenuItem) {
    this.menuItemService.setActiveMenuItem(menuItem);
  }

  trackByMenuItem(item: MenuItem): string {
    return item.url;
  }

  private loadMenuItems(): void {
    combineLatest([
      this.menuItemService.menuItems$,
      this.menuItemService.subMenuItems$,
    ])
      .pipe(
        tap(([menuItems, subMenuItems]) => {
          this.menuItems = menuItems;
          this.subMenuItems = subMenuItems;
          this.allMenuItems = [...menuItems, ...subMenuItems]; // Merge here
        }),
        takeUntil(this.destroy$)
      )
      .subscribe();
  }

  private updateActiveMenuItem(): void {
    if (!this.allMenuItems.length) {
      return;
    }
    const currentUrl = this.router.url.split('/').pop();

    this.menuItemService.setActiveMenuItem(
      this.allMenuItems.find((item) => item.url === '/' + currentUrl)
    );
  }

  private subscribeToRouterEvents(): void {
    this.router.events.pipe(take(1), takeUntil(this.destroy$)).subscribe(() => {
      this.updateActiveMenuItem();
    });
  }
}

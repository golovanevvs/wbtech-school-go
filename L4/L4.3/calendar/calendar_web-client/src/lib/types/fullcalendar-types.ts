// Интерфейсы для FullCalendar TypeScript типов

export interface Duration {
  years?: number;
  months?: number;
  days?: number;
  milliseconds?: number;
}

export interface EventObject {
  id?: string;
  title: string;
  start?: Date | string | null;
  end?: Date | string | null;
  allDay?: boolean;
  groupId?: string;
  backgroundColor?: string;
  borderColor?: string;
  textColor?: string;
  className?: string | string[];
  editable?: boolean;
  duration?: Duration;
  extendedProps?: { [key: string]: any };
}

export interface ViewObject {
  type: string;
  title: string;
  activeStart: Date;
  activeEnd: Date;
  currentStart: Date;
  currentEnd: Date;
}

export const DateUtils = {
  toDate(date: Date | string | null | undefined): Date | null {
    if (!date) return null;
    if (date instanceof Date) return date;
    return new Date(date);
  },
  
  toISOString(date: Date | string | null | undefined): string | null {
    const dateObj = this.toDate(date);
    return dateObj ? dateObj.toISOString() : null;
  },
  
  isDate(value: any): value is Date {
    return value instanceof Date;
  }
};

export interface EventResizeInfo {
  event: EventObject;
  relatedEvents?: EventObject[];
  oldEvent: EventObject;
  endDelta?: Duration;
  startDelta?: Duration;
  revert: () => void;
  view: ViewObject;
  el: HTMLElement;
  jsEvent: MouseEvent;
}

export interface EventResizeStartInfo {
  event: EventObject;
  jsEvent: MouseEvent;
  view: ViewObject;
}

export interface EventResizeStopInfo {
  event: EventObject;
  jsEvent: MouseEvent;
  view: ViewObject;
}

export interface EventDropInfo {
  event: EventObject;
  relatedEvents?: EventObject[];
  oldEvent: EventObject;
  delta: Duration;
  revert: () => void;
  view: ViewObject;
  el: HTMLElement;
  jsEvent: MouseEvent;
}

export interface EventClickInfo {
  event: EventObject;
  el: HTMLElement;
  jsEvent: MouseEvent;
  view: ViewObject;
}

export interface EventMouseEnterInfo {
  event: EventObject;
  el: HTMLElement;
  jsEvent: MouseEvent;
  view: ViewObject;
}

export interface EventMouseLeaveInfo {
  event: EventObject;
  el: HTMLElement;
  jsEvent: MouseEvent;
  view: ViewObject;
}

export interface FullCalendarEventResizeInfo {
  event: {
    id?: string;
    title: string;
    start?: Date | null;
    end?: Date | null;
    allDay?: boolean;
    [key: string]: any;
  };
  relatedEvents?: Array<{
    id?: string;
    title: string;
    start?: Date | null;
    end?: Date | null;
    allDay?: boolean;
    [key: string]: any;
  }>;
  oldEvent: {
    id?: string;
    title: string;
    start?: Date | null;
    end?: Date | null;
    allDay?: boolean;
    [key: string]: any;
  };
  endDelta?: Duration;
  startDelta?: Duration;
  revert: () => void;
  view: ViewObject;
  el: HTMLElement;
  jsEvent: MouseEvent;
}

export interface FullCalendarEventResizeStartInfo {
  event: {
    id?: string;
    title: string;
    start?: Date | null;
    end?: Date | null;
    allDay?: boolean;
    [key: string]: any;
  };
  jsEvent: MouseEvent;
  view: ViewObject;
}

export interface FullCalendarEventResizeStopInfo {
  event: {
    id?: string;
    title: string;
    start?: Date | null;
    end?: Date | null;
    allDay?: boolean;
    [key: string]: any;
  };
  jsEvent: MouseEvent;
  view: ViewObject;
}

export interface FullCalendarEventDropInfo {
  event: {
    id?: string;
    title: string;
    start?: Date | null;
    end?: Date | null;
    allDay?: boolean;
    [key: string]: any;
  };
  relatedEvents?: Array<{
    id?: string;
    title: string;
    start?: Date | null;
    end?: Date | null;
    allDay?: boolean;
    [key: string]: any;
  }>;
  oldEvent: {
    id?: string;
    title: string;
    start?: Date | null;
    end?: Date | null;
    allDay?: boolean;
    [key: string]: any;
  };
  delta: Duration;
  revert: () => void;
  view: ViewObject;
  el: HTMLElement;
  jsEvent: MouseEvent;
}
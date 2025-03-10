import { Component, Input } from '@angular/core';
import { NgFor } from '@angular/common';
import { AnsiPipe } from './ansi.pipe';

interface LogEntry {
  line: string;
  timestamp: string;
}

@Component({
  selector: 'app-table',
  standalone: true,
  imports: [NgFor, AnsiPipe],
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.css']
})
export class TableComponent {
  @Input() logs: LogEntry[] = [];

  formatTimestamp(timestamp: string): string {
    return new Date(timestamp).toLocaleString();
  }
} 
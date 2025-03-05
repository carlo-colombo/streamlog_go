import {Component, OnInit} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {SseClient} from 'ngx-sse-client';
import {HttpHeaders} from '@angular/common/http';
import {FilterComponent} from './filter/filter.component';
import {TableComponent} from './table/table.component';

interface LogEntry {
  line: string;
  timestamp: string;
}

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [FilterComponent, TableComponent],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  title = 'app';
  logs: LogEntry[] = [];

  constructor(
    private sseClient: SseClient,
  ) {
    const headers = new HttpHeaders().set('Authorization', `Basic YWRtaW46YWRtaW4=`);

    this.sseClient.stream('/logs?sse', {
      keepAlive: true,
      reconnectionDelay: 1_000,
      responseType: 'event'
    }, {headers}, 'GET')
      .subscribe((event) => {
        if (event.type === 'error') {
          const errorEvent = event as ErrorEvent;
          console.error(errorEvent.error, errorEvent.message);
        } else {
          const messageEvent = event as MessageEvent;
          
          if (messageEvent.type === 'reset') {
            this.logs = [];
          } else if (messageEvent.data) {
            const logEntry: LogEntry = JSON.parse(messageEvent.data);
            this.logs.unshift(logEntry);
          }
        }
      });
  }

  ngOnInit(): void {
    // throw new Error('Method not implemented.');
  }

  formatTimestamp(timestamp: string): string {
    return new Date(timestamp).toLocaleString();
  }
}

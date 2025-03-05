import {Component, OnInit} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {SseClient} from 'ngx-sse-client';
import {HttpHeaders} from '@angular/common/http';
import {NgFor} from '@angular/common';
import {DomSanitizer, SafeHtml} from '@angular/platform-browser';
import {FormsModule} from '@angular/forms';
import {HttpClient} from '@angular/common/http';

interface LogEntry {
  line: string;
  timestamp: string;
}

@Component({
  selector: 'app-root',
  imports: [NgFor, FormsModule],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css'
})
export class AppComponent implements OnInit {
  title = 'app';
  logs: LogEntry[] = [];
  filter: string = '';

  constructor(
    private sseClient: SseClient,
    private sanitizer: DomSanitizer,
    private http: HttpClient
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

  updateFilter() {
    this.http.post('/filter', { filter: this.filter }, {
      headers: new HttpHeaders().set('Authorization', `Basic YWRtaW46YWRtaW4=`)
    }).subscribe();
  }

  formatTimestamp(timestamp: string): string {
    return new Date(timestamp).toLocaleString();
  }
}
